package response

import (
	"ApkAdmin/model/project"
	"time"
)

// AnnouncementWithReadStatus 公告及阅读状态
type AnnouncementWithReadStatus struct {
	project.SystemAnnouncement
	IsRead   int        `json:"is_read" gorm:"column:is_read"`     // 是否已读
	ReadTime *time.Time `json:"read_time" gorm:"column:read_time"` // 阅读时间
	IsClosed int        `json:"is_closed" gorm:"column:is_closed"` // 是否关闭
}
