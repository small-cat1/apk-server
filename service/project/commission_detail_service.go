package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	"ApkAdmin/model/project/request"
	"ApkAdmin/model/project/response"
	"errors"
	"gorm.io/gorm"
	"time"
)

type CommissionDetailService struct{}

// GetCommissionDetailList 分页获取分佣明细列表
func (s *CommissionDetailService) GetCommissionDetailList(search request.CommissionDetailSearch) (list []project.CommissionDetail, total int64, err error) {
	limit := search.PageSize
	offset := search.PageSize * (search.Page - 1)
	db := global.GVA_DB.Model(&project.CommissionDetail{})
	// 用户ID筛选
	if search.UserId != nil && *search.UserId > 0 {
		db = db.Where("user_id = ?", *search.UserId)
	}
	// 状态筛选
	if search.Status != "" && search.Status != "all" {
		db = db.Where("status = ?", search.Status)
	}

	// 时间筛选
	if search.TimeFilter != "" && search.TimeFilter != "all" {
		db = s.applyTimeFilter(db, search.TimeFilter)
	}

	// 查询总数
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// 分页查询，按创建时间倒序
	err = db.Order("create_time DESC").Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}

// ClientGetCommissionDetailList GetCommissionDetailList 分页获取分佣明细列表（包含统计信息）
func (s *CommissionDetailService) ClientGetCommissionDetailList(userId uint, search request.ClientCommissionDetailSearch) (response response.CommissionDetailListResponse, err error) {
	limit := search.PageSize
	offset := search.PageSize * (search.Page - 1)

	// 构建查询条件
	db := global.GVA_DB.Model(&project.CommissionDetail{})

	// 状态筛选
	if search.Status != "" && search.Status != "all" {
		db = db.Where("status = ?", search.Status)
	}

	// 时间筛选
	if search.TimeFilter != "" && search.TimeFilter != "all" {
		db = s.applyTimeFilter(db, search.TimeFilter)
	}

	// 查询总数
	err = db.Count(&response.Total).Error
	if err != nil {
		return response, err
	}

	// 分页查询列表，按创建时间倒序
	err = db.Order("create_time DESC").Limit(limit).Offset(offset).Find(&response.List).Error
	if err != nil {
		return response, err
	}

	// 查询统计信息（使用相同的筛选条件）
	statsDb := global.GVA_DB.Model(&project.CommissionDetail{})
	clientSearch := request.CommissionDetailSearch{
		ClientCommissionDetailSearch: search,
	}
	statsDb = s.applyFilters(statsDb, clientSearch)

	err = statsDb.Select("COALESCE(SUM(commission), 0) as total_commission, COUNT(*) as total_orders").
		Scan(&response.Stats).Error
	if err != nil {
		return response, err
	}

	// 设置分页信息
	response.Page = search.Page
	response.PageSize = search.PageSize
	return response, nil
}

// GetCommissionDetailById 根据ID获取分佣明细
func (s *CommissionDetailService) GetCommissionDetailById(id uint) (detail project.CommissionDetail, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&detail).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return detail, errors.New("分佣明细不存在")
		}
		return detail, err
	}
	return detail, nil
}

// GetCommissionStats 获取佣金统计信息
func (s *CommissionDetailService) GetCommissionStats(userId uint, status string, timeFilter string) (stats response.CommissionStats, err error) {
	db := global.GVA_DB.Model(&project.CommissionDetail{})

	// 用户ID筛选
	if userId > 0 {
		db = db.Where("user_id = ?", userId)
	}

	// 状态筛选
	if status != "" && status != "all" {
		db = db.Where("status = ?", status)
	}

	// 时间筛选
	if timeFilter != "" && timeFilter != "all" {
		db = s.applyTimeFilter(db, timeFilter)
	}

	// 查询累计佣金和订单数
	err = db.Select("COALESCE(SUM(commission), 0) as total_commission, COUNT(*) as total_orders").
		Scan(&stats).Error

	return stats, err
}

// applyFilters 应用筛选条件（统一的筛选逻辑）
func (s *CommissionDetailService) applyFilters(db *gorm.DB, search request.CommissionDetailSearch) *gorm.DB {
	// 用户ID筛选
	if search.UserId != nil && *search.UserId > 0 {
		db = db.Where("user_id = ?", *search.UserId)
	}

	// 状态筛选
	if search.Status != "" && search.Status != "all" {
		db = db.Where("status = ?", search.Status)
	}

	// 时间筛选
	if search.TimeFilter != "" && search.TimeFilter != "all" {
		db = s.applyTimeFilter(db, search.TimeFilter)
	}

	return db
}
func (s *CommissionDetailService) applyTimeFilter(db *gorm.DB, timeFilter string) *gorm.DB {
	now := time.Now()

	switch timeFilter {
	case "today":
		// 今天：从今天0点到现在
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		db = db.Where("create_time >= ?", startOfDay)

	case "week":
		// 本周：从本周一0点到现在
		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7 // 将周日从0改为7
		}
		startOfWeek := now.AddDate(0, 0, -(weekday - 1))
		startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, now.Location())
		db = db.Where("create_time >= ?", startOfWeek)

	case "month":
		// 本月：从本月1号0点到现在
		startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		db = db.Where("create_time >= ?", startOfMonth)
	}

	return db
}
