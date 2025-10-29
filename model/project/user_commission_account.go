package project

import (
	"errors"
	"time"
)

type UserCommissionAccount struct {
	ID              int64     `gorm:"primarykey;comment:账户ID" json:"id"`
	UserID          uint      `gorm:"not null;uniqueIndex;comment:用户ID" json:"userId"`
	AvailableAmount float64   `gorm:"type:decimal(10,2);default:0.00;comment:可提现金额" json:"availableAmount"`
	FrozenAmount    float64   `gorm:"type:decimal(10,2);default:0.00;comment:冻结金额" json:"frozenAmount"`
	TotalEarnings   float64   `gorm:"type:decimal(10,2);default:0.00;comment:累计收益" json:"totalEarnings"`
	WithdrawnAmount float64   `gorm:"type:decimal(10,2);default:0.00;comment:已提现金额" json:"withdrawnAmount"`
	CreatedAt       time.Time `gorm:"comment:创建时间" json:"createdAt"`
	UpdatedAt       time.Time `gorm:"comment:更新时间" json:"updatedAt"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (UserCommissionAccount) TableName() string {
	return "user_commission_account"
}

// GetTotalAmount 获取总金额（可用+冻结）
func (a *UserCommissionAccount) GetTotalAmount() float64 {
	return a.AvailableAmount + a.FrozenAmount
}

// CanWithdraw 检查是否可以提现指定金额
func (a *UserCommissionAccount) CanWithdraw(amount float64) bool {
	return a.AvailableAmount >= amount && amount > 0
}

// AddEarnings 增加收益
func (a *UserCommissionAccount) AddEarnings(amount float64) {
	a.AvailableAmount += amount
	a.TotalEarnings += amount
}

// FreezeAmount 冻结金额
func (a *UserCommissionAccount) FreezeAmount(amount float64) error {
	if !a.CanWithdraw(amount) {
		return errors.New("可用余额不足")
	}
	a.AvailableAmount -= amount
	a.FrozenAmount += amount
	return nil
}

// UnfreezeAmount 解冻金额（提现失败时）
func (a *UserCommissionAccount) UnfreezeAmount(amount float64) {
	a.FrozenAmount -= amount
	a.AvailableAmount += amount
}

// CompleteWithdraw 完成提现（从冻结金额扣除）
func (a *UserCommissionAccount) CompleteWithdraw(amount float64) {
	a.FrozenAmount -= amount
	a.WithdrawnAmount += amount
}

// UserCommissionAccountSimple 创建一个只包含需要字段的结构体
type UserCommissionAccountSimple struct {
	ID              int64   `json:"id"`
	UserID          uint    `json:"user_id"`
	AvailableAmount float64 `json:"available_amount"`
	TotalEarnings   float64 `json:"total_earnings"`
}

func (UserCommissionAccountSimple) TableName() string {
	return "user_commission_account"
}
