package project

import "time"

// AppScreenshot 应用截图表
type AppScreenshot struct {
	ID             uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	AppID          string    `json:"app_id" gorm:"size:100;not null;comment:应用ID"`
	LanguageCode   string    `json:"language_code" gorm:"size:10;comment:语言代码"`
	Platform       string    `json:"platform" gorm:"type:enum('android','ios','harmony','windows');not null;comment:平台"`
	ScreenshotURL  string    `json:"screenshot_url" gorm:"size:500;not null;comment:截图URL"`
	ScreenshotType string    `json:"screenshot_type" gorm:"type:enum('phone','tablet','desktop','watch');default:'phone';comment:截图类型"`
	DisplayOrder   int       `json:"display_order" gorm:"default:0;comment:显示顺序"`
	CreatedAt      time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"comment:更新时间"`
}

// TableName 指定表名
func (AppScreenshot) TableName() string {
	return "app_screenshots"
}
