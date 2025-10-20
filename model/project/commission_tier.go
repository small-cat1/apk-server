package project

import "time"

// CommissionTier 阶梯分佣等级表模型
type CommissionTier struct {
	ID              int        `gorm:"primarykey;column:id" json:"id"`
	Name            string     `gorm:"column:name;type:varchar(50);not null;comment:等级名称：青铜推广员、白银推广员等" json:"name"`
	MinSubordinates int        `gorm:"column:min_subordinates;type:int;not null;default:0;comment:最低直属下级人数" json:"minSubordinates"`
	Rate            float64    `gorm:"column:rate;type:decimal(5,2);not null;comment:分佣比例(%)" json:"rate"`
	Color           string     `gorm:"column:color;type:varchar(20);comment:等级颜色标识" json:"color"`
	Icon            string     `gorm:"column:icon;type:varchar(100);comment:等级图标" json:"icon"`
	Sort            int        `gorm:"column:sort;type:int;default:0;comment:排序（数字越大等级越高）" json:"sort"`
	Status          *int       `gorm:"column:status;type:tinyint(1);default:1;comment:状态：1-启用, 0-禁用" json:"status"`
	CreateTime      *time.Time `gorm:"column:create_time;comment:创建时间" json:"createTime"`
	UpdateTime      *time.Time `gorm:"column:update_time;comment:更新时间" json:"updateTime"`
}

// TableName 指定表名
func (CommissionTier) TableName() string {
	return "commission_tiers"
}

// IsEnabled 判断等级是否启用
func (c *CommissionTier) IsEnabled() bool {
	return *c.Status == 1
}
