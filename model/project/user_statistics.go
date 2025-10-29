package project

import "time"

// UserStatistics 用户统计表
type UserStatistics struct {
	UserID              uint       `json:"user_id" gorm:"primaryKey;comment:用户ID"`
	TotalDownloads      uint       `json:"total_downloads" gorm:"default:0;comment:总下载次数"`
	TotalSpent          float64    `json:"total_spent" gorm:"type:decimal(10,2);default:0.00;comment:总消费金额"`
	TotalOrders         uint       `json:"total_orders" gorm:"default:0;comment:总订单数"`
	SuccessfulReferrals uint       `json:"successful_referrals" gorm:"default:0;comment:成功推荐人数"`
	LastDownloadAt      *time.Time `json:"last_download_at" gorm:"comment:最后下载时间"`
	LastOrderAt         *time.Time `json:"last_order_at" gorm:"comment:最后订单时间"`
	CreatedAt           time.Time  `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt           time.Time  `json:"updated_at" gorm:"comment:更新时间"`

	// 关联关系
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (UserStatistics) TableName() string {
	return "user_statistics"
}
