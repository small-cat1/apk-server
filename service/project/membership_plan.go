package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	"ApkAdmin/model/project/request"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type MembershipPlanService struct {
}

func (a *MembershipPlanService) Exists(code string) (bool, error) {
	var count int64
	err := global.GVA_DB.Model(&project.MembershipPlan{}).Where("plan_code = ?", code).Count(&count).Error
	return count > 0, err
}

func (a *MembershipPlanService) ExistsExcludeID(code string, excludeID uint) (bool, error) {
	var count int64
	err := global.GVA_DB.Model(&project.MembershipPlan{}).Where("plan_code = ? AND id != ?", code, excludeID).Count(&count).Error
	return count > 0, err
}

func (a *MembershipPlanService) GetMembershipPlan(id uint) (membershipPlan project.MembershipPlan, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&membershipPlan).Error
	return
}

func (a *MembershipPlanService) GetByID(id uint) (*project.MembershipPlan, error) {
	var membershipPlan project.MembershipPlan
	err := global.GVA_DB.First(&membershipPlan, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &membershipPlan, nil
}

func (a *MembershipPlanService) CreateMembershipPlan(req request.MembershipPlanCreateRequest) (err error) {
	// 检查会员套餐是否已存在
	exists, err := a.Exists(req.PlanCode)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("会员套餐 %s 已存在", req.PlanName)
	}
	membershipPlan := req.ToMembershipPlan()
	err = global.GVA_DB.Create(&membershipPlan).Error
	return err
}

func (a *MembershipPlanService) UpdateMembershipPlan(req *request.MembershipPlanUpdateRequest) (err error) {
	// 检查记录是否存在
	existing, err := a.GetByID(req.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("记录不存在")
	}
	// 检查会员套餐是否被其他记录使用
	if existing.PlanCode != req.PlanCode {
		exists, err := a.ExistsExcludeID(req.PlanCode, req.ID)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("会员套餐 %s 已被其他记录使用", req.PlanName)
		}
	}
	existing.PlanCode = req.PlanCode
	existing.PlanName = req.PlanName
	existing.PlanType = req.PlanType
	existing.Platform = req.Platform
	existing.DurationDays = req.DurationDays
	existing.BasePrice = &req.BasePrice
	existing.CurrencyCode = req.CurrencyCode
	existing.DiscountPercentage = &req.DiscountPercentage
	existing.FinalPrice = &req.FinalPrice
	existing.DownloadLimitDaily = req.DownloadLimitDaily
	existing.DownloadLimitMonthly = req.DownloadLimitMonthly
	existing.IsActive = req.IsActive
	existing.IsFeatured = req.IsFeatured
	existing.SortOrder = req.SortOrder
	existing.Description = req.Description
	err = global.GVA_DB.Debug().Omit("created_at").Updates(existing).Error
	return err
}

func (a *MembershipPlanService) DeleteMembershipPlan(cid int) (err error) {
	// 检查记录是否存在
	existing, err := a.GetByID(uint(cid))
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("记录不存在")
	}

	var total int64
	err = global.GVA_DB.Model(&project.UserMembership{}).Where("plan_id = ?", cid).Count(&total).Error
	if err != nil {
		return err
	}
	if total > 0 {
		return errors.New("该会员套餐下存在用户使用，无法删除")
	}

	var t1 int64
	err = global.GVA_DB.Model(&project.MembershipOrder{}).Where("plan_id = ?", cid).Count(&t1).Error
	if err != nil {
		return err
	}
	if t1 > 0 {
		return errors.New("该会员套餐下存在订单")
	}
	return global.GVA_DB.Model(&project.MembershipPlan{}).Where("id = ?", cid).Delete(&project.MembershipPlan{}).Error
}

func (a *MembershipPlanService) GetMembershipPlanList(info request.MembershipPlanListRequest, order string, desc bool) (list interface{}, total int64, err error) {
	var membershipPlanLists []project.MembershipPlan
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 构建查询条件
	db := global.GVA_DB.Model(&project.MembershipPlan{})
	db = a.buildSearchConditions(db, info)
	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return membershipPlanLists, total, err
	}
	// 分页和排序
	db = db.Limit(limit).Offset(offset)
	// 构建排序条件
	orderStr := a.buildOrderConditions(order, desc)
	err = db.Order(orderStr).Find(&membershipPlanLists).Error
	return membershipPlanLists, total, err
}

func (a *MembershipPlanService) GetAllMembershipPlan() (list interface{}, err error) {
	var plans []project.MembershipPlan
	err = global.GVA_DB.Model(&project.MembershipPlan{}).
		Where("is_active = ?", 1).
		Find(&plans).Error
	if err != nil {
		return nil, err
	}
	return plans, err
}

// buildSearchConditions 构建搜索条件
func (a *MembershipPlanService) buildSearchConditions(db *gorm.DB, info request.MembershipPlanListRequest) *gorm.DB {
	// 套餐代码搜索（支持模糊搜索）
	if info.PlanCode != "" {
		db = db.Where("plan_code LIKE ?", "%"+info.PlanCode+"%")
	}

	// 套餐名称搜索（模糊搜索）
	if info.PlanName != "" {
		db = db.Where("plan_name LIKE ?", "%"+info.PlanName+"%")
	}

	// 套餐类型搜索（精确匹配）
	if info.PlanType != "" {
		db = db.Where("plan_type = ?", info.PlanType)
	}

	// 平台搜索（JSON数组字段查询）
	if info.Platform != "" {
		// 使用JSON_CONTAINS函数查询JSON数组字段
		// MySQL: JSON_CONTAINS(platform, '"android"')
		// PostgreSQL: platform::jsonb ? 'android'
		// SQLite: JSON_EXTRACT(platform, '$[*]') LIKE '%android%'

		// 根据数据库类型选择合适的查询方式
		switch global.GVA_CONFIG.System.DbType {
		case "mysql":
			db = db.Where("JSON_CONTAINS(platform, ?)", `"`+info.Platform+`"`)
		case "postgres":
			db = db.Where("platform::jsonb ? ?", info.Platform)
		default: // sqlite等其他数据库
			db = db.Where("JSON_EXTRACT(platform, '$[*]') LIKE ?", "%"+info.Platform+"%")
		}
	}

	// 货币代码搜索
	if info.CurrencyCode != "" {
		db = db.Where("currency_code = ?", info.CurrencyCode)
	}

	// 状态搜索
	if info.IsActive != nil {
		db = db.Where("is_active = ?", *info.IsActive)
	}

	// 推荐状态搜索
	if info.IsFeatured != nil {
		db = db.Where("is_featured = ?", *info.IsFeatured)
	}

	// 价格范围搜索
	if info.MinPrice != nil && *info.MinPrice > 0 {
		db = db.Where("final_price >= ?", *info.MinPrice)
	}
	if info.MaxPrice != nil && *info.MaxPrice > 0 {
		db = db.Where("final_price <= ?", *info.MaxPrice)
	}

	// 创建时间范围搜索
	if info.StartDate != "" {
		db = db.Where("created_at >= ?", info.StartDate)
	}
	if info.EndDate != "" {
		db = db.Where("created_at <= ?", info.EndDate+" 23:59:59")
	}

	// 关键字搜索（同时搜索套餐代码、名称、描述）
	if info.Keyword != "" {
		keyword := "%" + info.Keyword + "%"
		db = db.Where("plan_code LIKE ? OR plan_name LIKE ? OR description LIKE ?",
			keyword, keyword, keyword)
	}

	return db
}

// buildOrderConditions 构建排序条件
func (a *MembershipPlanService) buildOrderConditions(order string, desc bool) string {
	// 默认排序
	defaultOrder := "sort_order ASC, created_at DESC"

	if order == "" {
		return defaultOrder
	}

	// 验证排序字段安全性
	allowedOrderFields := map[string]bool{
		"id":          true,
		"plan_code":   true,
		"plan_name":   true,
		"plan_type":   true,
		"final_price": true,
		"base_price":  true,
		"sort_order":  true,
		"created_at":  true,
		"updated_at":  true,
		"is_active":   true,
		"is_featured": true,
	}

	if !allowedOrderFields[order] {
		return defaultOrder
	}

	orderStr := order
	if desc {
		orderStr += " DESC"
	} else {
		orderStr += " ASC"
	}

	// 添加二级排序
	if order != "sort_order" {
		orderStr += ", sort_order ASC"
	}
	if order != "created_at" {
		orderStr += ", created_at DESC"
	}

	return orderStr
}
