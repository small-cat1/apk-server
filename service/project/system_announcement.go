package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	"ApkAdmin/model/project/request"
	"ApkAdmin/model/project/response"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type SystemAnnouncementService struct {
}

// MarkAsRead 标记公告为已读（使用 ON DUPLICATE KEY UPDATE）
func (s *SystemAnnouncementService) MarkAsRead(userID int64, announcementID int64) error {
	// 1. 检查公告是否存在
	var count int64
	if err := global.GVA_DB.Model(&project.SystemAnnouncement{}).
		Where("id = ? AND status = 1", announcementID).
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return errors.New("公告不存在或已下线")
	}

	// 2. 使用原生 SQL 执行 UPSERT
	now := time.Now()
	sql := `
        INSERT INTO user_announcement_reads 
            (user_id, announcement_id, is_read, read_time, created_at) 
        VALUES 
            (?, ?, 1, ?, ?)
        ON DUPLICATE KEY UPDATE 
            is_read = 1,
            read_time = ?
    `

	return global.GVA_DB.Exec(sql, userID, announcementID, now, now, now).Error
}

// GetAnnouncement 获取公告详情
func (s *SystemAnnouncementService) GetAnnouncement(conditions ...func(*gorm.DB) *gorm.DB) (res *project.SystemAnnouncement, err error) {
	query := global.GVA_DB.Model(&project.SystemAnnouncement{})
	// 应用所有条件
	for _, condition := range conditions {
		query = condition(query)
	}
	var announcement project.SystemAnnouncement
	err = query.First(&announcement).Error
	return &announcement, err
}

func (s SystemAnnouncementService) ListAnnouncement(info request.ListAnnouncementRequest) (list []project.SystemAnnouncement, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&project.SystemAnnouncement{})
	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	// 分页查询
	err = db.Order("created_at desc").Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}

func (s SystemAnnouncementService) PageInfoAnnouncement(info request.ListAnnouncementRequest, userID uint) (list []response.AnnouncementWithReadStatus, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	// 先统计总数
	countDB := global.GVA_DB.Model(&project.SystemAnnouncement{}).Where("status = ?", 1)
	if info.Type == 1 || info.Type == 2 || info.Type == 3 {
		countDB.Where("type = ?", info.Type)
	}
	err = countDB.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// LEFT JOIN 查询公告和阅读状态
	db := global.GVA_DB.
		Table("system_announcements sa").
		Select(`
            sa.*,
            COALESCE(uar.is_read, 0) as is_read,
            uar.read_time,
            COALESCE(uar.is_closed, 0) as is_closed
        `).
		Joins("LEFT JOIN user_announcement_reads uar ON sa.id = uar.announcement_id AND uar.user_id = ?", userID).
		Where("sa.status = ?", 1)

	if info.Type == 1 || info.Type == 2 || info.Type == 3 {
		db.Where("sa.type = ?", info.Type)
	}

	err = db.Order("sa.created_at desc").
		Limit(limit).
		Offset(offset).
		Scan(&list).Error

	return list, total, err
}

func (s SystemAnnouncementService) CreateAnnouncement(info request.CreateAnnouncementRequest, userID uint) error {
	announcement := &project.SystemAnnouncement{
		Title:       info.Title,
		Content:     info.Content,
		Type:        info.Type,
		DisplayType: info.DisplayType,
		TargetUsers: info.TargetUsers,
		LinkURL:     info.LinkUrl,
		IsClosable:  info.IsClosable,
		Status:      info.Status,
		CreatedBy:   userID,
	}
	// 通过 Get 方法获取时间（可能为 nil）
	if startTime := info.GetStartTime(); startTime != nil {
		announcement.StartTime = startTime
	}
	if endTime := info.GetEndTime(); endTime != nil {
		announcement.EndTime = endTime
	}
	return global.GVA_DB.Model(&project.SystemAnnouncement{}).Create(&announcement).Error
}

func (s SystemAnnouncementService) UpdateAnnouncement(info request.UpdateAnnouncementRequest) error {

	// 检查公告是否存在
	var existingAnnouncement project.SystemAnnouncement
	if err := global.GVA_DB.First(&existingAnnouncement, info.ID).Error; err != nil {
		return fmt.Errorf("公告不存在: %v", err)
	}

	// 构建更新数据
	updates := make(map[string]interface{})

	if info.Title != "" {
		updates["title"] = info.Title
	}
	if info.Content != "" {
		updates["content"] = info.Content
	}
	if info.Type != 0 {
		updates["type"] = info.Type
	}
	if info.DisplayType != 0 {
		updates["display_type"] = info.DisplayType
	}
	if info.TargetUsers != "" {
		updates["target_users"] = info.TargetUsers
	}
	if info.LinkUrl != "" {
		updates["link_url"] = info.LinkUrl
	}
	if info.IsClosable != 0 {
		updates["is_closable"] = info.IsClosable
	}
	if info.Status != nil {
		updates["status"] = info.Status
	}

	// 获取转换后的时间
	if startTime := info.GetStartTime(); startTime != nil {
		updates["start_time"] = startTime
	}
	if endTime := info.GetEndTime(); endTime != nil {
		updates["end_time"] = endTime
	}

	// 执行更新
	return global.GVA_DB.Model(&project.SystemAnnouncement{}).
		Where("id = ?", info.ID).
		Updates(updates).Error
}

func (s SystemAnnouncementService) DeleteAnnouncements(ids []uint) error {
	// 批量删除
	return global.GVA_DB.Where("id IN ?", ids).Delete(&project.SystemAnnouncement{}).Error
}
