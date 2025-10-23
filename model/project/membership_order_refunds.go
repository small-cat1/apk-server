package project

import (
	"ApkAdmin/global"
	"time"
)

// MembershipOrderRefund 会员订单退款记录
type MembershipOrderRefund struct {
	global.GVA_MODEL
	OrderID            uint       `json:"order_id" gorm:"not null;comment:订单ID;index"`
	OrderNo            string     `json:"order_no" gorm:"type:varchar(32);not null;comment:订单号;index"`
	RefundAmount       float64    `json:"refund_amount" gorm:"type:decimal(10,2);not null;comment:退款金额"`
	RefundReason       string     `json:"refund_reason" gorm:"type:text;comment:退款原因"`
	RefundType         string     `json:"refund_type" gorm:"type:enum('full','partial');default:full;comment:退款类型"`
	RefundStatus       string     `json:"refund_status" gorm:"type:enum('pending','processing','success','failed','cancelled');default:pending;comment:退款状态"`
	ThirdPartyRefundID string     `json:"third_party_refund_id" gorm:"type:varchar(100);comment:第三方退款ID"`
	OperatorID         *uint      `json:"operator_id" gorm:"comment:操作员ID"`
	OperatorName       string     `json:"operator_name" gorm:"type:varchar(50);comment:操作员姓名"`
	ProcessedAt        *time.Time `json:"processed_at" gorm:"comment:处理时间"`
	CompletedAt        *time.Time `json:"completed_at" gorm:"comment:完成时间"`
	FailureReason      string     `json:"failure_reason" gorm:"type:text;comment:失败原因"`
	Metadata           *string    `json:"metadata" gorm:"type:json;comment:额外数据"`

	// 关联查询
	Order    Order       `json:"order" gorm:"foreignKey:OrderID;references:ID"`
	Operator interface{} `json:"operator" gorm:"-"` // 这里可以关联到用户表
}

func (MembershipOrderRefund) TableName() string {
	return "membership_order_refunds"
}

// RefundStatusOptions 退款状态选项
var RefundStatusOptions = []string{
	"pending",    // 待处理
	"processing", // 处理中
	"success",    // 成功
	"failed",     // 失败
	"cancelled",  // 已取消
}

// RefundTypeOptions 退款类型选项
var RefundTypeOptions = []string{
	"full",    // 全额退款
	"partial", // 部分退款
}

// GetRefundStatusLabel 获取退款状态标签
func (r *MembershipOrderRefund) GetRefundStatusLabel() string {
	labels := map[string]string{
		"pending":    "待处理",
		"processing": "处理中",
		"success":    "退款成功",
		"failed":     "退款失败",
		"cancelled":  "已取消",
	}
	if label, exists := labels[r.RefundStatus]; exists {
		return label
	}
	return r.RefundStatus
}

// GetRefundTypeLabel 获取退款类型标签
func (r *MembershipOrderRefund) GetRefundTypeLabel() string {
	labels := map[string]string{
		"full":    "全额退款",
		"partial": "部分退款",
	}
	if label, exists := labels[r.RefundType]; exists {
		return label
	}
	return r.RefundType
}

// IsCompleted 是否已完成
func (r *MembershipOrderRefund) IsCompleted() bool {
	return r.RefundStatus == "success"
}

// IsFailed 是否失败
func (r *MembershipOrderRefund) IsFailed() bool {
	return r.RefundStatus == "failed"
}

// CanCancel 是否可以取消
func (r *MembershipOrderRefund) CanCancel() bool {
	return r.RefundStatus == "pending"
}

// CanRetry 是否可以重试
func (r *MembershipOrderRefund) CanRetry() bool {
	return r.RefundStatus == "failed"
}
