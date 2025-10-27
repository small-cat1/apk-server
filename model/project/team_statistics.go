package project

import "time"

type TeamStatistics struct {
	ID               int64   `gorm:"primarykey;comment:统计ID" json:"id"`
	UserID           int64   `gorm:"uniqueIndex;not null;comment:用户ID" json:"userId"`
	TotalMembers     int     `gorm:"default:0;comment:直属下级总人数" json:"totalMembers"`
	TodayNew         int     `gorm:"default:0;comment:今日新增直属下级" json:"todayNew"`
	ActiveMembers    int     `gorm:"default:0;comment:活跃直属下级数（近30天有消费）" json:"activeMembers"`
	TotalConsumption float64 `gorm:"type:decimal(12,2);default:0.00;comment:直属下级总消费" json:"totalConsumption"`
	TotalCommission  float64 `gorm:"type:decimal(10,2);default:0.00;comment:累计获得佣金" json:"totalCommission"`
	// 等级相关
	CurrentTierID *int      `gorm:"index;comment:当前阶梯等级ID" json:"currentTierId"`
	CreatedAt     time.Time `gorm:"comment:创建时间" json:"createdAt"`
	UpdatedAt     time.Time `gorm:"comment:更新时间" json:"updatedAt"`

	// 关联
	CurrentTier *CommissionTier `gorm:"foreignKey:CurrentTierID" json:"currentTier,omitempty"`
	User        *User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (TeamStatistics) TableName() string {
	return "team_statistics"
}
