package project

import "time"

// UserAnnouncementRead 用户公告阅读记录
type UserAnnouncementRead struct {
	ID             int64      `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	UserID         int64      `json:"user_id" gorm:"not null;comment:用户ID"`
	AnnouncementID int64      `json:"announcement_id" gorm:"not null;comment:公告ID"`
	IsRead         int        `json:"is_read" gorm:"default:0;comment:是否已读"`
	ReadTime       *time.Time `json:"read_time" gorm:"comment:阅读时间"`
	IsClosed       int        `json:"is_closed" gorm:"default:0;comment:是否已关闭横幅"`
	ClosedTime     *time.Time `json:"closed_time" gorm:"comment:关闭时间"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
}

// TableName 指定表名
func (UserAnnouncementRead) TableName() string {
	return "user_announcement_reads"
}
