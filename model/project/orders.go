package project

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// ==================== 订单表 ====================

// Order 订单
type Order struct {
	ID                   uint64             `gorm:"primarykey;autoIncrement" json:"id"`
	OrderNo              string             `gorm:"type:varchar(32);not null;uniqueIndex:uk_order_no;comment:订单号" json:"orderNo"`
	UserID               uint               `gorm:"not null;index:idx_user_id;comment:用户ID" json:"userId"`
	OrderType            OrderType          `gorm:"type:enum('membership','account_product');not null;comment:订单类型" json:"orderType"`
	ProductID            uint               `gorm:"not null;comment:商品ID" json:"productId"`
	ProductCode          string             `gorm:"type:varchar(50);not null;comment:商品代码快照" json:"productCode"`
	ProductName          string             `gorm:"type:varchar(100);not null;comment:商品名称快照" json:"productName"`
	MembershipSubType    *MembershipSubType `gorm:"type:enum('new','renew','upgrade','downgrade');comment:会员订单子类型" json:"membershipSubType,omitempty"`
	UpgradeCredit        float64            `gorm:"type:decimal(10,2);default:0.00;comment:升级抵扣金额" json:"upgradeCredit"`
	PreviousMembershipID *uint              `gorm:"comment:升级前的会员记录ID" json:"previousMembershipId,omitempty"`
	Quantity             uint               `gorm:"default:1;comment:购买数量" json:"quantity"`
	AccountIDs           AccountIDList      `gorm:"type:json;comment:分配的账号ID列表" json:"accountIds,omitempty"`
	OriginalPrice        float64            `gorm:"type:decimal(10,2);not null;comment:原价" json:"originalPrice"`
	DiscountAmount       float64            `gorm:"type:decimal(10,2);default:0.00;comment:优惠金额" json:"discountAmount"`
	FinalAmount          float64            `gorm:"type:decimal(10,2);not null;comment:最终金额" json:"finalAmount"`
	CurrencyCode         string             `gorm:"type:varchar(3);not null;default:CNY;comment:货币代码" json:"currencyCode"`
	PaymentMethod        *string            `gorm:"type:varchar(50);comment:支付方式" json:"paymentMethod,omitempty"`
	PaymentID            *string            `gorm:"type:varchar(100);comment:第三方支付ID" json:"paymentId,omitempty"`
	Status               OrderStatus        `gorm:"type:enum('pending','paid','failed','refunded','cancelled');not null;default:pending;index:idx_status;comment:订单状态" json:"status"`
	PaidAt               *time.Time         `gorm:"comment:支付时间" json:"paidAt,omitempty"`
	PaymentDeadline      *time.Time         `gorm:"index:idx_payment_deadline;comment:支付截止时间" json:"paymentDeadline,omitempty"`
	ExpiredAt            time.Time          `gorm:"comment:订单过期时间" json:"expiredAt,omitempty"`
	CreatedAt            time.Time          `gorm:"not null;comment:创建时间" json:"createdAt"`
	UpdatedAt            time.Time          `gorm:"not null;comment:更新时间" json:"updatedAt"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (Order) TableName() string {
	return "orders"
}

// ==================== 自定义类型 ====================

// OrderType 订单类型
type OrderType string

const (
	OrderTypeMembership     OrderType = "membership"      // 会员订单
	OrderTypeAccountProduct OrderType = "account_product" // 账号产品订单
)

// MembershipSubType 会员订单子类型
type MembershipSubType string

const (
	MembershipSubTypeNew       MembershipSubType = "new"       // 新购
	MembershipSubTypeRenew     MembershipSubType = "renew"     // 续费
	MembershipSubTypeUpgrade   MembershipSubType = "upgrade"   // 升级
	MembershipSubTypeDowngrade MembershipSubType = "downgrade" // 降级
)

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // 待支付
	OrderStatusPaid      OrderStatus = "paid"      // 已支付
	OrderStatusFailed    OrderStatus = "failed"    // 支付失败
	OrderStatusRefunded  OrderStatus = "refunded"  // 已退款
	OrderStatusCancelled OrderStatus = "cancelled" // 已取消
)

// AccountIDList 账号ID列表（用于JSON字段）
type AccountIDList []uint

// Scan 实现 sql.Scanner 接口
func (a *AccountIDList) Scan(value interface{}) error {
	if value == nil {
		*a = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, a)
}

// Value 实现 driver.Valuer 接口
func (a AccountIDList) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	return json.Marshal(a)
}

// ==================== Order 辅助方法 ====================

// IsPending 是否待支付
func (o *Order) IsPending() bool {
	return o.Status == OrderStatusPending
}

// IsPaid 是否已支付
func (o *Order) IsPaid() bool {
	return o.Status == OrderStatusPaid
}

// CanPay 是否可以支付
func (o *Order) CanPay() bool {
	return o.Status == OrderStatusPending &&
		(o.PaymentDeadline == nil || time.Now().Before(*o.PaymentDeadline))
}

// CanCancel 是否可以取消
func (o *Order) CanCancel() bool {
	return o.Status == OrderStatusPending
}

// CanRefund 是否可以退款
func (o *Order) CanRefund() bool {
	return o.Status == OrderStatusPaid
}

// IsExpired 是否已过期
func (o *Order) IsExpired() bool {
	return time.Now().After(o.ExpiredAt)
}

// IsMembershipOrder 是否是会员订单
func (o *Order) IsMembershipOrder() bool {
	return o.OrderType == OrderTypeMembership
}

// IsAccountProductOrder 是否是账号产品订单
func (o *Order) IsAccountProductOrder() bool {
	return o.OrderType == OrderTypeAccountProduct
}
