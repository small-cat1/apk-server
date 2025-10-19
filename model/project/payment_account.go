package project

import (
	"ApkAdmin/global"
	"encoding/json"
	"gorm.io/gorm"
	"strings"
	"time"
)

// PaymentAccount 支付账号配置表
type PaymentAccount struct {
	global.GVA_MODEL
	Name         string `json:"name" gorm:"type:varchar(100);not null;comment:账号名称"`
	ProviderCode string `json:"provider_code" gorm:"type:varchar(50);not null;index;comment:支付服务商代码"`
	AccountType  string `json:"account_type" gorm:"type:enum('personal','enterprise');default:personal;comment:账号类型"`

	// 通用配置 - JSON存储所有配置参数
	Config string `json:"config" gorm:"type:json;not null;comment:支付配置参数"`

	// 账号状态和权重
	Status         string  `json:"status" gorm:"type:enum('active','inactive','maintenance');default:active;comment:状态"`
	Weight         int     `json:"weight" gorm:"default:1;comment:权重(用于负载均衡)"`
	MaxDailyAmount float64 `json:"max_daily_amount" gorm:"type:decimal(15,2);default:0;comment:日限额(0表示无限制)"`

	// 使用统计
	DailyAmount float64    `json:"daily_amount" gorm:"type:decimal(15,2);default:0;comment:当日交易金额"`
	TotalAmount float64    `json:"total_amount" gorm:"type:decimal(15,2);default:0;comment:总交易金额"`
	TotalOrders int64      `json:"total_orders" gorm:"default:0;comment:总订单数"`
	LastUsedAt  *time.Time `json:"last_used_at" gorm:"comment:最后使用时间"`

	// 分组和标签
	Group  string `json:"group" gorm:"type:varchar(50);comment:分组"`
	Tags   string `json:"tags" gorm:"type:varchar(255);comment:标签"`
	Region string `json:"region" gorm:"type:varchar(50);comment:服务地区"`
	Remark string `json:"remark" gorm:"type:text;comment:备注"`

	// 关联查询
	Provider PaymentProvider `json:"provider" gorm:"foreignKey:ProviderCode;references:Code"`
}

func (PaymentAccount) TableName() string {
	return "payment_accounts"
}

// BeforeCreate 创建前钩子
func (pa *PaymentAccount) BeforeCreate(tx *gorm.DB) error {
	// 可以在这里设置默认值或执行其他逻辑
	if pa.Status == "" {
		pa.Status = "active"
	}
	if pa.Weight == 0 {
		pa.Weight = 1
	}
	return nil
}

// BeforeUpdate 更新前钩子
func (pa *PaymentAccount) BeforeUpdate(tx *gorm.DB) error {
	// 可以在这里执行更新前的验证或其他逻辑
	return nil
}

// AfterCreate 创建后钩子
func (pa *PaymentAccount) AfterCreate(tx *gorm.DB) error {
	// 可以在这里记录操作日志或发送通知
	return nil
}

// AfterUpdate 更新后钩子
func (pa *PaymentAccount) AfterUpdate(tx *gorm.DB) error {
	// 可以在这里记录操作日志或发送通知
	return nil
}

// AfterDelete 删除后钩子
func (pa *PaymentAccount) AfterDelete(tx *gorm.DB) error {
	// 可以在这里清理相关数据或记录日志
	return nil
}

// GetConfigMap 获取配置参数的Map格式
func (pa *PaymentAccount) GetConfigMap() (map[string]interface{}, error) {
	var config map[string]interface{}
	if pa.Config == "" {
		return config, nil
	}
	err := json.Unmarshal([]byte(pa.Config), &config)
	return config, err
}

// SetConfigFromMap 从Map设置配置参数
func (pa *PaymentAccount) SetConfigFromMap(config map[string]interface{}) error {
	configBytes, err := json.Marshal(config)
	if err != nil {
		return err
	}
	pa.Config = string(configBytes)
	return nil
}

// IsActive 判断账号是否可用
func (pa *PaymentAccount) IsActive() bool {
	return pa.Status == "active"
}

// IsDailyLimitReached 判断是否达到日限额
func (pa *PaymentAccount) IsDailyLimitReached() bool {
	if pa.MaxDailyAmount <= 0 {
		return false
	}
	return pa.DailyAmount >= pa.MaxDailyAmount
}

// GetDailyUsageRate 获取日限额使用率
func (pa *PaymentAccount) GetDailyUsageRate() float64 {
	if pa.MaxDailyAmount <= 0 {
		return 0
	}
	return pa.DailyAmount / pa.MaxDailyAmount * 100
}

// CanProcess 判断是否可以处理交易
func (pa *PaymentAccount) CanProcess(amount float64) bool {
	if !pa.IsActive() {
		return false
	}

	if pa.MaxDailyAmount > 0 && (pa.DailyAmount+amount) > pa.MaxDailyAmount {
		return false
	}

	return true
}

// AddTransaction 添加交易记录（更新统计数据）
func (pa *PaymentAccount) AddTransaction(amount float64) {
	pa.DailyAmount += amount
	pa.TotalAmount += amount
	pa.TotalOrders++
	now := time.Now()
	pa.LastUsedAt = &now
}

// GetTagsSlice 获取标签切片
func (pa *PaymentAccount) GetTagsSlice() []string {
	if pa.Tags == "" {
		return []string{}
	}

	tags := strings.Split(pa.Tags, ",")
	var result []string
	for _, tag := range tags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			result = append(result, tag)
		}
	}
	return result
}

// SetTagsFromSlice 从切片设置标签
func (pa *PaymentAccount) SetTagsFromSlice(tags []string) {
	pa.Tags = strings.Join(tags, ",")
}

// AddTags 添加标签
func (pa *PaymentAccount) AddTags(newTags []string) {
	existingTags := pa.GetTagsSlice()
	tagMap := make(map[string]bool)

	// 添加现有标签
	for _, tag := range existingTags {
		tagMap[tag] = true
	}

	// 添加新标签
	for _, tag := range newTags {
		tag = strings.TrimSpace(tag)
		if tag != "" {
			tagMap[tag] = true
		}
	}

	// 转换为切片
	var tags []string
	for tag := range tagMap {
		tags = append(tags, tag)
	}

	pa.SetTagsFromSlice(tags)
}

// RemoveTags 移除标签
func (pa *PaymentAccount) RemoveTags(removeTags []string) {
	existingTags := pa.GetTagsSlice()
	removeMap := make(map[string]bool)

	// 构建移除标签的映射
	for _, tag := range removeTags {
		removeMap[strings.TrimSpace(tag)] = true
	}

	// 过滤标签
	var tags []string
	for _, tag := range existingTags {
		if !removeMap[tag] {
			tags = append(tags, tag)
		}
	}

	pa.SetTagsFromSlice(tags)
}
