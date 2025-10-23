package project

import "time"

// AccountFlow 流水表：记录历史明细
type AccountFlow struct {
	ID            int64   `gorm:"primarykey" json:"id"`
	UserID        int64   `gorm:"not null;index" json:"userId"`
	Type          string  `gorm:"type:varchar(20);not null;index" json:"type"`
	Amount        float64 `gorm:"type:decimal(10,2);not null" json:"amount"`
	BalanceBefore float64 `gorm:"type:decimal(10,2);not null" json:"balanceBefore"`
	BalanceAfter  float64 `gorm:"type:decimal(10,2);not null" json:"balanceAfter"`

	// ✅ 独立的关联字段
	OrderID    *int64 `gorm:"index" json:"orderId,omitempty"`
	WithdrawID *int64 `gorm:"index" json:"withdrawId,omitempty"`
	RefundID   *int64 `gorm:"index" json:"refundId,omitempty"`

	FlowNo     string    `gorm:"type:varchar(32);not null;uniqueIndex" json:"flowNo"`
	Remark     string    `gorm:"type:varchar(255)" json:"remark"`
	CreateTime time.Time `gorm:"index" json:"createTime"`

	// 关联（可选）
	Order    *Order          `gorm:"foreignKey:OrderID" json:"order,omitempty"`
	Withdraw *WithdrawRecord `gorm:"foreignKey:WithdrawID" json:"withdraw,omitempty"`
}

// 流水类型常量
const (
	FlowTypeCommissionIn = "commission_in" // 佣金收入
	FlowTypeWithdrawOut  = "withdraw_out"  // 提现支出
	FlowTypeFreeze       = "freeze"        // 冻结
	FlowTypeUnfreeze     = "unfreeze"      // 解冻
	FlowTypeRefund       = "refund"        // 退款
)
