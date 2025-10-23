package project

import "time"

// ==================== 提现记录表 ====================

// WithdrawRecord 提现记录
type WithdrawRecord struct {
	ID           int64      `gorm:"primarykey;comment:提现ID" json:"id"`
	UserID       int64      `gorm:"not null;index:idx_user_id;comment:用户ID" json:"userId"`
	WithdrawNo   string     `gorm:"type:varchar(32);not null;uniqueIndex;comment:提现单号" json:"withdrawNo"`
	Amount       float64    `gorm:"type:decimal(10,2);not null;comment:提现金额" json:"amount"`
	Fee          float64    `gorm:"type:decimal(10,2);default:0.00;comment:手续费" json:"fee"`
	ActualAmount float64    `gorm:"type:decimal(10,2);not null;comment:实际到账金额" json:"actualAmount"`
	WithdrawType string     `gorm:"type:varchar(20);not null;comment:提现方式：alipay-支付宝, wechat-微信" json:"withdrawType"`
	AccountName  *string    `gorm:"type:varchar(50);comment:账户名" json:"accountName,omitempty"`
	AccountNo    *string    `gorm:"type:varchar(100);comment:账户号" json:"accountNo,omitempty"`
	Status       string     `gorm:"type:varchar(20);default:pending;index:idx_status;comment:状态" json:"status"`
	RejectReason *string    `gorm:"type:varchar(255);comment:拒绝原因" json:"rejectReason,omitempty"`
	AuditTime    *time.Time `gorm:"comment:审核时间" json:"auditTime,omitempty"`
	CompleteTime *time.Time `gorm:"comment:完成时间" json:"completeTime,omitempty"`
	Remark       *string    `gorm:"type:varchar(255);comment:备注" json:"remark,omitempty"`
	CreateTime   time.Time  `gorm:"index:idx_create_time;comment:创建时间" json:"createTime"`
	UpdateTime   time.Time  `gorm:"comment:更新时间" json:"updateTime"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (WithdrawRecord) TableName() string {
	return "withdraw_records"
}

// 提现方式常量
const (
	WithdrawTypeAlipay = "alipay" // 支付宝
	WithdrawTypeWechat = "wechat" // 微信
)

// 提现状态常量
const (
	WithdrawStatusPending   = "pending"   // 待审核
	WithdrawStatusApproved  = "approved"  // 已通过
	WithdrawStatusRejected  = "rejected"  // 已拒绝
	WithdrawStatusCompleted = "completed" // 已完成
)

// IsPending 是否待审核
func (w *WithdrawRecord) IsPending() bool {
	return w.Status == WithdrawStatusPending
}

// IsCompleted 是否已完成
func (w *WithdrawRecord) IsCompleted() bool {
	return w.Status == WithdrawStatusCompleted
}

// CanCancel 是否可以取消
func (w *WithdrawRecord) CanCancel() bool {
	return w.Status == WithdrawStatusPending
}
