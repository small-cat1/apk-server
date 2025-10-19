package project

import (
	"ApkAdmin/constants"
	"encoding/json"
	"time"
)

// MembershipPlan 会员套餐结构体
type MembershipPlan struct {
	ID                   uint               `json:"id" gorm:"primaryKey;autoIncrement;comment:套餐ID"`
	PlanCode             string             `json:"plan_code" gorm:"type:varchar(50);uniqueIndex:uk_plan_code;not null;comment:套餐代码"`
	PlanName             string             `json:"plan_name" gorm:"type:varchar(100);not null;comment:套餐名称"`
	PlanType             constants.PlanType `json:"plan_type" gorm:"type:enum('monthly','yearly','lifetime');not null;index:idx_plan_type;comment:套餐类型"`
	Platform             json.RawMessage    `json:"platform" gorm:"type:enum('android','ios','harmony','windows');not null;comment:平台类型"`
	DurationDays         *int               `json:"duration_days" gorm:"comment:有效天数（终身会员为NULL）"`
	BasePrice            *float64           `json:"base_price" gorm:"type:decimal(10,2);not null;comment:基础价格"`
	CurrencyCode         string             `json:"currency_code" gorm:"type:varchar(3);not null;default:USD;comment:货币代码"`
	DiscountPercentage   *float64           `json:"discount_percentage" gorm:"type:decimal(5,2);default:0.00;comment:折扣百分比"`
	FinalPrice           *float64           `json:"final_price" gorm:"type:decimal(10,2);not null;comment:最终价格"`
	DownloadLimitDaily   *int               `json:"download_limit_daily" gorm:"comment:每日下载限制（NULL表示无限制）"`
	DownloadLimitMonthly *int               `json:"download_limit_monthly" gorm:"comment:每月下载限制"`
	IsActive             *bool              `json:"is_active" gorm:"type:tinyint(1);default:1;index:idx_is_active;comment:是否启用"`
	IsFeatured           *bool              `json:"is_featured" gorm:"type:tinyint(1);default:0;comment:是否推荐套餐"`
	SortOrder            *int               `json:"sort_order" gorm:"default:0;index:idx_sort_order;comment:排序权重"`
	Description          *string            `json:"description" gorm:"type:text;comment:套餐描述"`
	CreatedAt            time.Time          `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt            time.Time          `json:"updated_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP"`
}

// TableName 指定表名
func (MembershipPlan) TableName() string {
	return "membership_plans"
}

// IsLifetime 判断是否为终身套餐
func (m *MembershipPlan) IsLifetime() bool {
	return m.PlanType == constants.PlanTypeLifetime
}

// HasDownloadLimit 判断是否有下载限制
func (m *MembershipPlan) HasDownloadLimit() bool {
	return m.DownloadLimitDaily != nil || m.DownloadLimitMonthly != nil
}

// GetActualPrice 获取实际价格（考虑折扣）
func (m *MembershipPlan) GetActualPrice() float64 {
	if *m.DiscountPercentage > 0 {
		return *m.BasePrice * (1 - *m.DiscountPercentage/100)
	}
	return *m.BasePrice
}
