package project

import (
	"ApkAdmin/constants"
	"time"
)

type DownloadLog struct {
	ID           uint               `gorm:"primarykey" json:"id"`
	UserID       uint               `gorm:"not null;default:0;index:idx_user_app;index:idx_user_created" json:"userId" comment:"用户ID"`
	AppID        uint               `gorm:"not null;index:idx_user_app;index:idx_app_created" json:"appId" comment:"应用ID"`
	PackageID    *uint              `gorm:"index" json:"packageId" comment:"安装包ID"`
	MembershipID *uint              `gorm:"index" json:"membershipId" comment:"会员ID"`
	Platform     constants.Platform `gorm:"type:varchar(20);not null" json:"platform" comment:"平台"`
	Success      bool               `gorm:"not null;default:0" json:"success" comment:"是否成功"`
	FailReason   *string            `gorm:"type:varchar(255)" json:"failReason" comment:"失败原因"`
	IP           string             `gorm:"type:varchar(45);not null" json:"ip" comment:"IP地址"`
	UserAgent    *string            `gorm:"type:varchar(500)" json:"userAgent" comment:"用户代理"`
	DeviceType   *string            `gorm:"type:varchar(20)" json:"deviceType" comment:"设备类型"`
	CreatedAt    time.Time          `gorm:"not null;default:CURRENT_TIMESTAMP;index:idx_created;index:idx_user_created;index:idx_app_created" json:"createdAt"`
}

func (DownloadLog) TableName() string {
	return "download_logs"
}

// SetFailReason 设置失败原因
func (d *DownloadLog) SetFailReason(reason string) {
	d.FailReason = &reason
}

// SetUserAgent 设置用户代理
func (d *DownloadLog) SetUserAgent(ua string) {
	if ua != "" {
		// 截断过长的 UA
		if len(ua) > 500 {
			ua = ua[:500]
		}
		d.UserAgent = &ua
	}
}

// SetDeviceType 设置设备类型
func (d *DownloadLog) SetDeviceType(deviceType string) {
	if deviceType != "" {
		d.DeviceType = &deviceType
	}
}

// ==================== 统计查询结构体 ====================

// DownloadStats 下载统计
type DownloadStats struct {
	TotalDownloads int64   `json:"totalDownloads"`
	SuccessCount   int64   `json:"successCount"`
	FailCount      int64   `json:"failCount"`
	SuccessRate    float64 `json:"successRate"`
}

// AppDownloadStats 应用下载统计
type AppDownloadStats struct {
	AppID            uint   `json:"appId"`
	AppName          string `json:"appName"`
	TotalDownloads   int    `json:"totalDownloads"`
	IOSDownloads     int    `json:"iosDownloads"`
	AndroidDownloads int    `json:"androidDownloads"`
}

// UserDownloadStats 用户下载统计
type UserDownloadStats struct {
	UserID         uint       `json:"userId"`
	TodayCount     int64      `json:"todayCount"`
	MonthCount     int64      `json:"monthCount"`
	TotalCount     int64      `json:"totalCount"`
	LastDownloadAt *time.Time `json:"lastDownloadAt"`
}
