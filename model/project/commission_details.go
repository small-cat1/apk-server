package project

import (
	"time"
)

// CommissionDetail 分佣明细表
type CommissionDetail struct {
	ID             uint       `json:"id" gorm:"primarykey;comment:明细ID"`
	UserId         uint       `json:"userId" gorm:"not null;comment:获得佣金的用户ID（推广人）;index:idx_user_id"`
	OrderId        uint       `json:"orderId" gorm:"not null;comment:订单ID;index:idx_order_id"`
	OrderNo        string     `json:"orderNo" gorm:"type:varchar(32);not null;comment:订单号;index:idx_order_no"`
	OrderUserId    uint       `json:"orderUserId" gorm:"not null;comment:下单用户ID（直属下级）;index:idx_order_user_id"`
	OrderUsername  string     `json:"orderUsername" gorm:"type:varchar(50);comment:下单用户名"`
	OrderAmount    float64    `json:"orderAmount" gorm:"type:decimal(10,2);not null;comment:订单金额"`
	CommissionRate float64    `json:"commissionRate" gorm:"type:decimal(5,4);not null;comment:佣金比例(小数形式，如0.1表示10%)"`
	Commission     float64    `json:"commission" gorm:"type:decimal(10,2);not null;comment:佣金金额"`
	TierId         *int       `json:"tierId" gorm:"comment:阶梯等级ID;index:idx_tier_id"`
	TierName       string     `json:"tierName" gorm:"type:varchar(50);comment:阶梯等级名称（冗余字段，方便查询）"`
	Status         string     `json:"status" gorm:"type:varchar(20);default:pending;comment:状态：pending-待结算, settled-已结算, frozen-冻结;index:idx_status"`
	SettleTime     *time.Time `json:"settleTime" gorm:"comment:结算时间"`
	Remark         string     `json:"remark" gorm:"type:varchar(255);comment:备注"`
	CreateTime     time.Time  `json:"createTime" gorm:"autoCreateTime;comment:创建时间;index:idx_create_time"`
}

// TableName 表名
func (CommissionDetail) TableName() string {
	return "commission_details"
}
