package project

import (
	"ApkAdmin/constants"
	"ApkAdmin/global"
	"encoding/json"
	"time"
)

// OrderType 订单类型
type OrderType string

const (
	OrderTypeNew       OrderType = "new"       // 新购
	OrderTypeRenew     OrderType = "renew"     // 续费
	OrderTypeUpgrade   OrderType = "upgrade"   // 升级
	OrderTypeDowngrade OrderType = "downgrade" // 降级
)

// MembershipOrder 会员订单表
type MembershipOrder struct {
	global.GVA_MODEL
	OrderNo              string                `json:"order_no" gorm:"type:varchar(32);not null;uniqueIndex:uk_order_no;comment:订单号"`
	UserID               uint                  `json:"user_id" gorm:"not null;index:idx_user_id;comment:用户ID"`
	PlanID               uint                  `json:"plan_id" gorm:"not null;index:idx_plan_id;comment:套餐ID"`
	PlanCode             string                `json:"plan_code" gorm:"type:varchar(50);not null;comment:套餐代码快照"`
	PlanName             string                `json:"plan_name" gorm:"type:varchar(100);not null;comment:套餐名称快照"`
	PlanType             constants.PlanType    `json:"plan_type" gorm:"type:enum('monthly','yearly','lifetime');not null;comment:套餐类型"`
	Platform             json.RawMessage       `json:"platform" gorm:"type:enum('android','ios','harmony','windows');not null;comment:购买平台"`
	OrderType            OrderType             `json:"order_type" gorm:"type:enum('new','renew','upgrade','downgrade');not null;comment:订单类型"`
	OriginalPrice        float64               `json:"original_price" gorm:"type:decimal(10,2);not null;comment:原价"`
	DiscountAmount       float64               `json:"discount_amount" gorm:"type:decimal(10,2);default:0.00;comment:折扣金额"`
	UpgradeCredit        float64               `json:"upgrade_credit" gorm:"type:decimal(10,2);default:0.00;comment:升级抵扣金额（原套餐剩余价值）"`
	FinalAmount          float64               `json:"final_amount" gorm:"type:decimal(10,2);not null;comment:实际支付金额"`
	CurrencyCode         string                `json:"currency_code" gorm:"type:varchar(3);not null;default:USD;comment:货币代码"`
	PaymentMethod        *string               `json:"payment_method" gorm:"type:varchar(50);comment:支付方式"`
	PaymentID            *string               `json:"payment_id" gorm:"type:varchar(100);comment:第三方支付ID"`
	Status               constants.OrderStatus `json:"status" gorm:"type:enum('pending','paid','failed','refunded','cancelled');not null;default:pending;index:idx_status;comment:订单状态"`
	PaidAt               *time.Time            `json:"paid_at" gorm:"comment:支付时间"`
	ExpiresAt            time.Time             `json:"expires_at" gorm:"not null;comment:订单过期时间"`
	PreviousMembershipID *uint                 `json:"previous_membership_id" gorm:"index:idx_previous_membership;comment:升级前的会员记录ID"`
	Metadata             json.RawMessage       `json:"metadata" gorm:"type:json;comment:额外数据（支付凭证等）"`
}

// TableName 指定表名
func (MembershipOrder) TableName() string {
	return "membership_orders"
}
