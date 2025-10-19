package project

import (
	"gorm.io/gorm"
	"time"
)

// PaymentProvider 支付服务商表
type PaymentProvider struct {
	ID          uint           `gorm:"primarykey" json:"id"`                                            // 主键ID
	CreatedAt   time.Time      `json:"created_at"`                                                      // 创建时间
	UpdatedAt   time.Time      `json:"updated_at"`                                                      // 更新时间
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`                                                  // 删除时间
	Code        string         `json:"code" gorm:"type:varchar(50);uniqueIndex;not null;comment:服务商代码"` // wechat, alipay, stripe, paypal
	Name        string         `json:"name" gorm:"type:varchar(100);not null;comment:服务商名称"`
	Description string         `json:"description" gorm:"type:text;comment:描述"`
	Status      string         `json:"status" gorm:"type:enum('active','inactive');default:active;comment:状态"`
	Icon        string         `json:"icon" gorm:"type:varchar(255);comment:图标URL"`
	SortOrder   int            `json:"sort_order" gorm:"default:0;comment:排序"`
}

func (PaymentProvider) TableName() string {
	return "payment_providers"
}
