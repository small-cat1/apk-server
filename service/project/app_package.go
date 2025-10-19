package project

import (
	"ApkAdmin/constants"
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	"ApkAdmin/model/project/request"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type AppPackageService struct{}

func (a *AppPackageService) Exists(code string) (bool, error) {
	var count int64
	err := global.GVA_DB.Model(&project.AppPackage{}).Where("app_id = ?", code).Count(&count).Error
	return count > 0, err
}

func (a *AppPackageService) GetAppPackageManual(id uint) (appPackage project.AppPackage, err error) {
	// 首先获取 AppPackage 基础信息
	err = global.GVA_DB.Where("id = ?", id).
		Preload("Application").
		First(&appPackage).Error
	if err != nil {
		return
	}

	// 手动查询关联的套餐
	var planRelations []project.PackagePlanRelation
	err = global.GVA_DB.Where("app_package_id = ?", id).
		Preload("Package").
		Find(&planRelations).Error
	if err != nil {
		return
	}
	// 提取套餐信息
	var membershipPlans []project.MembershipPlan
	for _, relation := range planRelations {
		membershipPlans = append(membershipPlans, relation.Package)
	}
	appPackage.MembershipPlans = membershipPlans
	appPackage.PlanRelations = planRelations
	return
}

func (a *AppPackageService) GetByID(id uint) (*project.AppPackage, error) {
	var country project.AppPackage
	err := global.GVA_DB.First(&country, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &country, nil
}

func (a *AppPackageService) CreateAppPackage(uid uint, req request.AppPackageCreateRequest) (err error) {
	// 检查应用是否存在
	var app project.Application
	err = global.GVA_DB.Model(&project.Application{}).Where("app_id = ?", req.AppID).First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("应用不存在: %s", req.AppID)
		}
		return err
	}
	// 开始事务
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	// 创建安装包
	appPackage := req.ToAppPackage()
	appPackage.CountryCode = app.CountryCode
	appPackage.AppName = app.AppName
	appPackage.CreatedBy = uint64(uid)
	// 创建安装包记录
	if err := tx.Create(&appPackage).Error; err != nil {
		tx.Rollback()
		return err
	}
	// 如果有套餐列表，创建关联关系
	if len(req.PlanList) > 0 {
		// 验证套餐ID是否存在
		var planCount int64
		err = tx.Model(&project.MembershipPlan{}).Where("id IN ?", req.PlanList).Count(&planCount).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		if planCount != int64(len(req.PlanList)) {
			tx.Rollback()
			return fmt.Errorf("部分套餐ID不存在")
		}
		// 创建包套餐关联关系
		for _, planID := range req.PlanList {
			relation := project.PackagePlanRelation{
				AppPackageID:     appPackage.ID,
				MembershipPlanID: uint(planID),
			}
			if err := tx.Create(&relation).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	// 提交事务
	return tx.Commit().Error
}

func (a *AppPackageService) UpdateAppPackage(useID uint, req *request.AppPackageUpdateRequest) (err error) {
	// 检查记录是否存在
	existing, err := a.GetByID(req.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("记录不存在")
	}

	// 开始事务
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	// 更新安装包基本信息
	updates := map[string]interface{}{
		"platform":     req.Platform,
		"version_name": req.VersionName,
		"version_code": req.VersionCode,
		"file_url":     req.FileURL,
		"status":       req.Status,
		"updated_at":   time.Now(),
		"updated_by":   useID,
	}

	err = tx.Model(&project.AppPackage{}).Where("id = ?", req.ID).Updates(updates).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// 处理套餐关系
	if req.PlanList != nil { // 如果传入了套餐列表，则更新关系
		// 删除现有的套餐关系
		err = tx.Where("app_package_id = ?", req.ID).Delete(&project.PackagePlanRelation{}).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		// 如果有新的套餐列表，创建新的关联关系
		if len(req.PlanList) > 0 {
			// 验证套餐ID是否存在
			var planCount int64
			err = tx.Model(&project.MembershipPlan{}).Where("id IN ?", req.PlanList).Count(&planCount).Error
			if err != nil {
				tx.Rollback()
				return err
			}

			if planCount != int64(len(req.PlanList)) {
				tx.Rollback()
				return fmt.Errorf("部分套餐ID不存在")
			}

			// 创建新的关联关系
			for _, planID := range req.PlanList {
				relation := project.PackagePlanRelation{
					AppPackageID:     uint64(req.ID),
					MembershipPlanID: uint(planID),
				}
				if err := tx.Create(&relation).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}
	}

	// 提交事务
	return tx.Commit().Error
}

// DeleteAppPackage 删除安装包
func (a *AppPackageService) DeleteAppPackage(cid int) (err error) {
	// 检查记录是否存在
	existing, err := a.GetByID(uint(cid))
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("记录不存在")
	}

	// 开始事务
	tx := global.GVA_DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if tx.Error != nil {
		return tx.Error
	}

	// 删除套餐关联关系
	err = tx.Where("app_package_id = ?", cid).Delete(&project.PackagePlanRelation{}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("删除套餐关联关系失败: %v", err)
	}

	// 删除安装包记录
	err = tx.Where("id = ?", cid).Delete(&project.AppPackage{}).Error
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("删除安装包失败: %v", err)
	}

	// 提交事务
	return tx.Commit().Error
}

func (a *AppPackageService) BatchUpdateApkStatus(useID uint, ids []uint, status constants.PackageStatus) error {
	// 更新安装包基本信息
	updates := map[string]interface{}{
		"status":       status,
		"published_at": time.Now(),
		"published_by": useID,
	}
	err := global.GVA_DB.Model(&project.AppPackage{}).Where("id IN ?", ids).Updates(updates).Error
	return err
}

func (a *AppPackageService) GetAppPackageList(info request.AppPackageListRequest, order string, desc bool) (list interface{}, total int64, err error) {
	var countryLists []project.AppPackage
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&project.AppPackage{})
	if info.AppID != "" {
		db = db.Where("app_id = ?", info.AppID)
	}
	if info.Platform != "" {
		db = db.Where("platform = ?", info.Platform)
	}
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}
	if info.CountryCode != "" {
		db = db.Where("country_code = ?", info.CountryCode)
	}
	err = db.Count(&total).Error
	if err != nil {
		return countryLists, total, err
	}
	db = db.Limit(limit).Offset(offset)
	OrderStr := "id desc"
	if order != "" {
		OrderStr = order
		if desc {
			OrderStr = order + " desc"
		}
	}
	// 预加载应用信息和套餐关系
	err = db.Preload("Application", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, app_id, app_name, app_icon, description,country_code, status")
	}).Order(OrderStr).Find(&countryLists).Error
	return countryLists, total, err
}
