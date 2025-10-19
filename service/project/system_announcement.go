package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	"ApkAdmin/model/project/request"
	"fmt"
	"gorm.io/gorm"
)

type SystemAnnouncementService struct {
}

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

func (s SystemAnnouncementService) PageInfoAnnouncement(info request.ListAnnouncementRequest) (list []project.SystemAnnouncement, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Debug().Model(&project.SystemAnnouncement{}).Where("status = ?", 1)
	if info.Type == 1 || info.Type == 2 || info.Type == 3 {
		db.Where("type = ?", info.Type)
	}
	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	// 分页查询
	err = db.Order("created_at desc").Limit(limit).Offset(offset).Find(&list).Error
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
