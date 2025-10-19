package project

import "time"

type SystemAnnouncement struct {
	ID          int64      `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	Title       string     `json:"title" gorm:"type:varchar(200);not null;comment:公告标题"`
	Content     string     `json:"content" gorm:"type:text;comment:公告内容（支持富文本）"`
	Type        int        `json:"type" gorm:"not null;comment:类型：1=紧急，2=重要，3=普通"`
	DisplayType int        `json:"display_type" gorm:"comment:展示类型：1=横幅，2=弹窗，3=卡片，4=仅消息中心"`
	TargetUsers string     `json:"target_users" gorm:"type:varchar(50);comment:目标用户：all=全部，vip=VIP，new=新用户"`
	LinkURL     string     `json:"link_url" gorm:"type:varchar(500);column:link_url;comment:跳转链接"`
	StartTime   *time.Time `json:"start_time" gorm:"comment:开始显示时间"`
	EndTime     *time.Time `json:"end_time" gorm:"comment:结束显示时间"`
	IsClosable  int        `json:"is_closable" gorm:"default:1;comment:是否可关闭：0=否，1=是"`
	Status      *int       `json:"status" gorm:"default:1;comment:状态：0=草稿，1=发布"`
	CreatedBy   uint       `json:"created_by" gorm:"comment:创建人"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
}

// TableName 指定表名
func (SystemAnnouncement) TableName() string {
	return "system_announcements"
}
