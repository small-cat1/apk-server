package project

import "time"

// SystemConfig 系统配置表
type SystemConfig struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	Scope     string     `gorm:"type:varchar(255);not null;uniqueIndex:uk_scope_key;comment:作用域,website(站点配置),pay(支付)" json:"scope"`
	Name      string     `gorm:"type:varchar(255);not null;comment:配置名称" json:"name"`
	Key       string     `gorm:"type:varchar(255);not null;uniqueIndex:uk_scope_key;comment:配置key" json:"key"`
	Value     string     `gorm:"type:text;comment:配置值" json:"value"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// TableName 指定表名
func (SystemConfig) TableName() string {
	return "system_configs"
}
