package project

import (
	"ApkAdmin/constants"
	"ApkAdmin/global"
	projectModel "ApkAdmin/model/project"
	"time"
)

// DownloadLogService 下载日志服务
type DownloadLogService struct{}

// ==================== 记录日志 ====================

// CreateDownloadLog 创建下载日志
func (s *DownloadLogService) CreateDownloadLog(log *projectModel.DownloadLog) error {
	return global.GVA_DB.Create(log).Error
}

// CreateSuccessLog 记录成功日志
func (s *DownloadLogService) CreateSuccessLog(userID, appID uint, platform constants.Platform, ip, userAgent string) error {
	log := &projectModel.DownloadLog{
		UserID:    userID,
		AppID:     appID,
		Platform:  platform,
		Success:   true,
		IP:        ip,
		CreatedAt: time.Now(),
	}
	log.SetUserAgent(userAgent)
	return s.CreateDownloadLog(log)
}

// CreateFailLog 记录失败日志
func (s *DownloadLogService) CreateFailLog(userID, appID uint, platform constants.Platform, ip, userAgent, reason string) error {
	log := &projectModel.DownloadLog{
		UserID:    userID,
		AppID:     appID,
		Platform:  platform,
		Success:   false,
		IP:        ip,
		CreatedAt: time.Now(),
	}
	log.SetUserAgent(userAgent)
	log.SetFailReason(reason)
	return s.CreateDownloadLog(log)
}

// ==================== 查询日志 ====================

// GetUserDownloadLogs 获取用户下载日志
func (s *DownloadLogService) GetUserDownloadLogs(userID uint, limit int) ([]projectModel.DownloadLog, error) {
	var logs []projectModel.DownloadLog
	err := global.GVA_DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// GetAppDownloadLogs 获取应用下载日志
func (s *DownloadLogService) GetAppDownloadLogs(appID uint, limit int) ([]projectModel.DownloadLog, error) {
	var logs []projectModel.DownloadLog
	err := global.GVA_DB.
		Where("app_id = ?", appID).
		Order("created_at DESC").
		Limit(limit).
		Find(&logs).Error
	return logs, err
}

// ==================== 统计查询 ====================

// GetUserDownloadStats 获取用户下载统计
func (s *DownloadLogService) GetUserDownloadStats(userID uint) (*projectModel.UserDownloadStats, error) {
	stats := &projectModel.UserDownloadStats{
		UserID: userID,
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// 今日下载
	global.GVA_DB.Model(&projectModel.DownloadLog{}).
		Where("user_id = ? AND created_at >= ?", userID, today).
		Count(&stats.TodayCount)

	// 本月下载
	global.GVA_DB.Model(&projectModel.DownloadLog{}).
		Where("user_id = ? AND created_at >= ?", userID, monthStart).
		Count(&stats.MonthCount)

	// 总下载
	global.GVA_DB.Model(&projectModel.DownloadLog{}).
		Where("user_id = ?", userID).
		Count(&stats.TotalCount)

	// 最后下载时间
	var lastLog projectModel.DownloadLog
	err := global.GVA_DB.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		First(&lastLog).Error

	if err == nil {
		stats.LastDownloadAt = &lastLog.CreatedAt
	}

	return stats, nil
}

// GetAppDownloadStats 获取应用下载统计
func (s *DownloadLogService) GetAppDownloadStats(appID uint) (*projectModel.DownloadStats, error) {
	stats := &projectModel.DownloadStats{}

	// 总下载数
	err := global.GVA_DB.Model(&projectModel.DownloadLog{}).
		Where("app_id = ?", appID).
		Count(&stats.TotalDownloads).Error

	if err != nil {
		return nil, err
	}

	// 成功数
	global.GVA_DB.Model(&projectModel.DownloadLog{}).
		Where("app_id = ? AND success = ?", appID, true).
		Count(&stats.SuccessCount)

	// 失败数
	stats.FailCount = stats.TotalDownloads - stats.SuccessCount

	// 成功率
	if stats.TotalDownloads > 0 {
		stats.SuccessRate = float64(stats.SuccessCount) / float64(stats.TotalDownloads) * 100
	}

	return stats, nil
}

// GetDailyDownloadCount 获取某天的下载次数
func (s *DownloadLogService) GetDailyDownloadCount(userID uint, date time.Time) (int, error) {
	var count int64

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	err := global.GVA_DB.Model(&projectModel.DownloadLog{}).
		Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, startOfDay, endOfDay).
		Count(&count).Error

	return int(count), err
}

// GetMonthlyDownloadCount 获取某月的下载次数
func (s *DownloadLogService) GetMonthlyDownloadCount(userID uint, year int, month time.Month) (int, error) {
	var count int64

	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	endOfMonth := startOfMonth.AddDate(0, 1, 0)

	err := global.GVA_DB.Model(&projectModel.DownloadLog{}).
		Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, startOfMonth, endOfMonth).
		Count(&count).Error

	return int(count), err
}

// ==================== 清理日志 ====================

// DeleteOldLogs 删除旧日志
func (s *DownloadLogService) DeleteOldLogs(days int) error {
	cutoffDate := time.Now().AddDate(0, 0, -days)
	return global.GVA_DB.
		Where("created_at < ?", cutoffDate).
		Delete(&projectModel.DownloadLog{}).Error
}

// ==================== 高级统计 ====================

// GetTopDownloadApps 获取下载最多的应用
func (s *DownloadLogService) GetTopDownloadApps(limit int) ([]projectModel.AppDownloadStats, error) {
	var stats []projectModel.AppDownloadStats

	err := global.GVA_DB.Raw(`
		SELECT 
			app_id,
			COUNT(*) as total_downloads,
			SUM(CASE WHEN platform = 'ios' THEN 1 ELSE 0 END) as ios_downloads,
			SUM(CASE WHEN platform = 'android' THEN 1 ELSE 0 END) as android_downloads
		FROM download_logs
		GROUP BY app_id
		ORDER BY total_downloads DESC
		LIMIT ?
	`, limit).Scan(&stats).Error

	return stats, err
}

// GetDownloadTrend 获取下载趋势（最近N天）
func (s *DownloadLogService) GetDownloadTrend(days int) (map[string]int, error) {
	startDate := time.Now().AddDate(0, 0, -days)

	var results []struct {
		Date  string `gorm:"column:date"`
		Count int    `gorm:"column:count"`
	}

	err := global.GVA_DB.Raw(`
		SELECT 
			DATE(created_at) as date,
			COUNT(*) as count
		FROM download_logs
		WHERE created_at >= ?
		GROUP BY DATE(created_at)
		ORDER BY date
	`, startDate).Scan(&results).Error

	if err != nil {
		return nil, err
	}

	trend := make(map[string]int)
	for _, r := range results {
		trend[r.Date] = r.Count
	}

	return trend, nil
}
