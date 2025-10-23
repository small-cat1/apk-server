package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	"ApkAdmin/utils"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type WithdrawService struct {
}

func (s *WithdrawService) CreateWithdraw(userID int64, amount float64) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 1. 创建提现申请
		withdraw := project.WithdrawRecord{
			UserID:     userID,
			Amount:     amount,
			Status:     project.WithdrawStatusPending,
			CreateTime: time.Now(),
		}
		if err := tx.Create(&withdraw).Error; err != nil {
			return err
		}

		// 2. 锁定账户
		var account project.UserCommissionAccount
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).
			First(&account).Error; err != nil {
			return err
		}

		// 3. 检查余额
		if account.AvailableAmount < amount {
			return errors.New("余额不足")
		}

		balanceBefore := account.AvailableAmount
		balanceAfter := balanceBefore - amount

		// 4. 冻结金额
		if err := tx.Model(&account).Updates(map[string]interface{}{
			"available_amount": balanceAfter,
			"frozen_amount":    gorm.Expr("frozen_amount + ?", amount),
		}).Error; err != nil {
			return err
		}

		// 5. 创建冻结流水
		flowNo := utils.GenerateFlowNo("FRZ")
		freezeFlow := project.AccountFlow{
			UserID:        userID,
			Type:          project.FlowTypeFreeze,
			Amount:        amount,
			BalanceBefore: balanceBefore,
			BalanceAfter:  balanceAfter,
			OrderID:       nil,
			WithdrawID:    &withdraw.ID, // ✅ 明确关联提现记录
			RefundID:      nil,
			FlowNo:        flowNo,
			Remark:        "提现冻结",
			CreateTime:    time.Now(),
		}

		return tx.Create(&freezeFlow).Error
	})
}

// GetWithdrawFlows ✅ 查询某个提现的所有流水（冻结、解冻、支出）
func (s *WithdrawService) GetWithdrawFlows(withdrawID int64) ([]project.AccountFlow, error) {
	var flows []project.AccountFlow
	err := global.GVA_DB.Where("withdraw_id = ?", withdrawID).
		Order("create_time ASC").
		Find(&flows).Error
	return flows, err
}
