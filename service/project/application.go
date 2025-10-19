package project

import (
	"ApkAdmin/constants"
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	"ApkAdmin/model/project/request"
	"ApkAdmin/model/project/response"
	"ApkAdmin/utils"
	"ApkAdmin/utils/upload"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"mime/multipart"
)

type ApplicationService struct{}

func (a *ApplicationService) SearchApps(req request.SearchAppRequest) (list interface{}, total int64, err error) {
	var applicationList []project.Application
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	// 构建查询条件
	db := global.GVA_DB.Model(&project.Application{}).Where("app_name LIKE ?", "%"+req.Keyword+"%")
	// 分页和排序
	db = db.Limit(limit).Offset(offset)
	sortStr := "sort_order desc"
	if req.SortType == "hot" {
		sortStr = "is_hot desc"
	}
	if req.SortType == "new" {
		sortStr = "created_at desc"
	}
	if req.SortType == "rating" {
		sortStr = "rating desc"
	}
	if req.SortType == "downloads" {
		sortStr = "download_count desc"
	}
	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return applicationList, total, err
	}
	// 构建排序条件
	err = db.Order(sortStr).Find(&applicationList).Error
	if err != nil {
		return nil, 0, err
	}
	return applicationList, total, err
}

// FilterAppsByCategory 按照分类分页获取应用
func (a *ApplicationService) FilterAppsByCategory(req request.FilterAppRequest) (list interface{}, total int64, err error) {
	var cid []uint
	if req.CategoryId == 0 {
		// 查找账号分类
		err = global.GVA_DB.Model(&project.AppCategory{}).
			Where("account_status = ?", 1).
			Pluck("id", &cid).Error
	} else {
		var category project.AppCategory
		// 查找账号分类
		err = global.GVA_DB.Model(&project.AppCategory{}).
			Where("id = ? and account_status = ?", req.CategoryId, 1).
			First(&category).Error
		if err != nil {
			return nil, 0, err
		}
		cid = append(cid, req.CategoryId)
	}

	var applicationList []project.Application
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	// 构建查询条件
	db := global.GVA_DB.Model(&project.Application{}).Where("category_id in ?", cid)
	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	sortStr := "sort_order desc"
	if req.SortType == "hot" {
		sortStr = "is_hot desc"
	}
	if req.SortType == "new" {
		sortStr = "created_at desc"
	}
	if req.SortType == "rating" {
		sortStr = "rating desc"
	}
	if req.SortType == "downloads" {
		sortStr = "download_count desc"
	}
	// 构建排序条件
	err = db.Limit(limit).Offset(offset).Order(sortStr).Find(&applicationList).Error
	if err != nil {
		return nil, 0, err
	}
	return applicationList, total, err
}

// GetAppByAccountCategory   根据账号分类获取应用列表
func (a *ApplicationService) GetAppByAccountCategory(req request.FilterAccountAppRequest) (result []response.AccountAppResp, total int64, err error) {
	var cid []uint
	if req.CategoryId == 0 {
		// 查找账号分类
		err = global.GVA_DB.Model(&project.AppCategory{}).
			Where("account_status = ?", 1).
			Pluck("id", &cid).Error
	} else {
		var category project.AppCategory
		// 查找账号分类
		err = global.GVA_DB.Model(&project.AppCategory{}).
			Where("id = ? and account_status = ?", req.CategoryId, 1).
			First(&category).Error
		if err != nil {
			return nil, 0, err
		}
		cid = append(cid, req.CategoryId)
	}
	if len(cid) <= 0 {
		return nil, 0, errors.New("暂无可用的应用账号分类")
	}
	var applicationList []project.Application
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)
	// 构建查询条件
	db := global.GVA_DB.Model(&project.Application{}).
		Where("category_id in ? and status = ?", cid, constants.ApplicationStatusActive)
	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	// 构建排序条件
	sortStr := a.buildSortString(req.SortType)
	// 分页和排序
	err = db.Preload("Accounts").Order(sortStr).
		Limit(limit).
		Offset(offset).
		Find(&applicationList).Error
	if err != nil {
		return nil, 0, err
	}
	result = make([]response.AccountAppResp, 0, len(applicationList))
	for _, v := range applicationList {
		accountSalesAccount := v.AccountSalesAccount
		if accountSalesAccount == 0 {
			accountSalesAccount = utils.RandNumber(50, 999)
		}
		tmp := response.AccountAppResp{
			ID:           v.ID,
			AppID:        v.AppID,
			AppName:      v.AppName,
			AppIcon:      v.AppIcon,
			CategoryID:   v.CategoryID,
			AccountPrice: v.AccountPrice,
			Rating:       v.Rating,
			Stock:        len(v.Accounts),
			SalesCount:   accountSalesAccount,
		}
		result = append(result, tmp)
	}
	return result, total, nil
}

// buildSortString 根据排序类型构建排序字符串
func (a *ApplicationService) buildSortString(sortType string) string {
	switch sortType {
	case "price_asc":
		// 价格从低到高
		return "account_price asc, sort_order desc"
	case "price_desc":
		// 价格从高到低
		return "account_price desc, sort_order desc"
	case "time_desc":
		// 最新上架
		return "created_at desc, sort_order desc"
	case "time_asc":
		// 最早上架
		return "created_at asc, sort_order desc"
	case "sales_desc":
		// 购买量从高到低
		return "sales_count desc, sort_order desc"
	case "sales_asc":
		// 购买量从低到高
		return "sales_count asc, sort_order desc"
	case "default":
		fallthrough
	default:
		// 综合排序（默认按 sort_order 降序）
		return "sort_order desc, created_at desc"
	}
}

// GetHotOrRecommendApp 获取热门和推荐的应用
func (a *ApplicationService) GetHotOrRecommendApp() (list interface{}, err error) {
	var app []project.Application
	err = global.GVA_DB.Model(&project.Application{}).Where("is_hot = ? or is_recommend = ?", 1, 1).Scan(&app).Error
	if err != nil {
		return nil, err
	}
	return app, err
}

// Exists 检查应用ID是否存在
func (a *ApplicationService) Exists(appID string) (bool, error) {
	var count int64
	err := global.GVA_DB.Model(&project.Application{}).Where("app_id = ?", appID).Count(&count).Error
	return count > 0, err
}

// ExistsExcludeID 检查应用ID是否存在（排除指定ID）
func (a *ApplicationService) ExistsExcludeID(appID string, excludeID uint) (bool, error) {
	var count int64
	err := global.GVA_DB.Model(&project.Application{}).Where("app_id = ? AND id != ?", appID, excludeID).Count(&count).Error
	return count > 0, err
}

// GetApplication 根据ID获取应用详情
func (a *ApplicationService) GetApplication(id uint) (application project.Application, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&application).Error
	return
}

// GetByID 根据ID获取应用
func (a *ApplicationService) GetByID(id uint) (*project.Application, error) {
	var application project.Application
	err := global.GVA_DB.First(&application, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("应用不存在！")
		}
		return nil, err
	}
	return &application, nil
}

// GetAppDetail 获取APP详情
func (a *ApplicationService) GetAppDetail(id uint, platform constants.Platform) (app project.Application, err error) {
	query := global.GVA_DB.Model(&project.Application{}).Where("id = ?", id)
	err = query.Preload("Packages", func(db *gorm.DB) *gorm.DB {
		return db.Where("platform = ?", platform).Order("created_at DESC")
	}).
		Preload("Packages.MembershipPlans").
		First(&app).Error
	return app, err
}

// CreateApplication 创建应用
func (a *ApplicationService) CreateApplication(userID uint, req request.ApplicationCreateRequest) (err error) {
	var category project.AppCategory
	err = global.GVA_DB.Model(&project.AppCategory{}).Where("id = ?", req.CategoryID).Preload("Parent").First(&category).Error
	if err != nil {
		return errors.New("分类不存在")
	}
	application := req.ToApplication()
	if category.Parent != nil {
		application.CategoryID = &category.Parent.ID
		application.SubcategoryID = &category.ID
	} else {
		application.CategoryID = &category.ID
		application.SubcategoryID = &category.ParentID
	}
	randRating := utils.GenerateRatingWithDistribution()
	application.Rating = &randRating
	application.AppID = uuid.New().String()
	application.Status = constants.ApplicationStatusActive
	application.CreatedBy = int64(userID)
	err = global.GVA_DB.Create(&application).Error
	return err
}

// UpdateApplication 更新应用
func (a *ApplicationService) UpdateApplication(req *request.ApplicationUpdateRequest) (err error) {
	// 检查记录是否存在
	existing, err := a.GetByID(req.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("应用不存在")
	}
	var category project.AppCategory
	err = global.GVA_DB.Model(&project.AppCategory{}).Where("id = ?", req.CategoryID).Preload("Parent").First(&category).Error
	if err != nil {
		return errors.New("分类不存在")
	}
	if category.Parent != nil {
		fmt.Println(category.Parent.ID, category.ID)
		existing.CategoryID = &category.Parent.ID
		existing.SubcategoryID = &category.ID
	} else {
		existing.CategoryID = &category.ID
		existing.SubcategoryID = &category.ParentID
	}
	// 更新字段
	existing.AccountPrice = req.AccountPrice
	existing.AppName = req.AppName
	existing.AppIcon = req.AppIcon
	existing.Description = req.Description
	existing.IsHot = req.IsHot
	existing.IsRecommend = req.IsRecommend
	existing.IsFree = &req.IsFree
	existing.SortOrder = req.SortOrder
	err = global.GVA_DB.Omit("created_at").Updates(existing).Error
	return err
}

// DeleteApplication 删除应用
func (a *ApplicationService) DeleteApplication(cid int) (err error) {
	// 检查记录是否存在
	existing, err := a.GetByID(uint(cid))
	if err != nil {
		return err
	}
	if existing == nil {
		return fmt.Errorf("应用不存在")
	}

	// 检查是否有关联的版本记录
	var versionCount int64
	err = global.GVA_DB.Model(&project.AppPackage{}).Where("app_id = ?", existing.AppID).Count(&versionCount).Error
	if err != nil {
		return err
	}
	if versionCount > 0 {
		return errors.New("该应用下存在版本记录，无法删除")
	}

	// 检查是否有关联的应用截图
	var downloadCount int64
	err = global.GVA_DB.Model(&project.AppScreenshot{}).Where("app_id = ?", existing.AppID).Count(&downloadCount).Error
	if err != nil {
		return err
	}
	if downloadCount > 0 {
		return errors.New("该应用下存在应用截图，无法删除")
	}

	return global.GVA_DB.Model(&project.Application{}).Where("id = ?", cid).Delete(&project.Application{}).Error
}

// BatchDeleteApplications 批量删除应用
func (a *ApplicationService) BatchDeleteApplications(ids []uint) error {
	// 检查每个应用是否可以删除
	for _, id := range ids {
		if err := a.DeleteApplication(int(id)); err != nil {
			return fmt.Errorf("删除应用ID %d 失败: %v", id, err)
		}
	}
	return nil
}

// BatchUpdateApplicationStatus 批量更新应用状态
func (a *ApplicationService) BatchUpdateApplicationStatus(ids []uint, status constants.ApplicationStatus) error {
	err := global.GVA_DB.Model(&project.Application{}).Where("id IN ?", ids).Update("status", status).Error
	return err
}

// CloneApplication 克隆应用
func (a *ApplicationService) CloneApplication(uid uint, req request.ApplicationCloneRequest) (*project.Application, error) {
	// 获取源应用
	sourceApp, err := a.GetByID(req.SourceID)
	if err != nil {
		return nil, err
	}
	if sourceApp == nil {
		return nil, fmt.Errorf("源应用不存在")
	}

	// 创建克隆应用
	cloneApp := project.Application{
		AppID:         sourceApp.AppID + "_copy",
		AppName:       sourceApp.AppName,
		CategoryID:    sourceApp.CategoryID,
		SubcategoryID: sourceApp.SubcategoryID,
		AppIcon:       sourceApp.AppIcon,
		Description:   sourceApp.Description,
		Status:        constants.ApplicationStatusActive,
		CreatedBy:     int64(uid),
	}

	if err := global.GVA_DB.Create(&cloneApp).Error; err != nil {
		return nil, err
	}

	return &cloneApp, nil
}

// UploadApplicationIcon 上传应用图标
func (a *ApplicationService) UploadApplicationIcon(appID uint64, file *multipart.FileHeader) (string, error) {
	// 验证应用是否存在
	application, err := a.GetByID(uint(appID))
	if err != nil {
		return "", err
	}
	if application == nil {
		return "", fmt.Errorf("应用不存在")
	}

	// 保存文件
	//iconURL, err := utils.SaveUploadedFile(file, "application_icons")
	oss := upload.NewOss()
	iconURL, _, uploadErr := oss.UploadFile(file)
	if uploadErr != nil {
		return "", uploadErr
	}
	if err != nil {
		return "", err
	}

	// 更新应用图标URL
	if err := global.GVA_DB.Model(application).Update("app_icon", iconURL).Error; err != nil {
		return "", err
	}

	return iconURL, nil
}

// GetApplicationList 获取应用列表
func (a *ApplicationService) GetApplicationList(info request.ApplicationListRequest, order string, desc bool) (list interface{}, total int64, err error) {
	var applicationList []project.Application
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	// 构建查询条件
	db := global.GVA_DB.Model(&project.Application{})
	db = a.buildSearchConditions(db, info)

	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return applicationList, total, err
	}

	// 分页和排序
	db = db.Limit(limit).Offset(offset)

	// 构建排序条件
	orderStr := a.buildOrderConditions(order, desc)
	err = db.Order(orderStr).Find(&applicationList).Error
	return applicationList, total, err
}

// buildSearchConditions 构建搜索条件
func (a *ApplicationService) buildSearchConditions(db *gorm.DB, info request.ApplicationListRequest) *gorm.DB {
	// 应用名称搜索（模糊搜索）
	if info.AppName != "" {
		db = db.Where("app_name LIKE ?", "%"+info.AppName+"%")
	}

	// 应用ID搜索（模糊搜索）
	if info.AppID != "" {
		db = db.Where("app_id LIKE ?", "%"+info.AppID+"%")
	}

	// 分类搜索
	if info.CategoryID > 0 {
		db = db.Where("category_id = ?", info.CategoryID)
	}

	// 状态搜索
	if info.Status != "" {
		db = db.Where("status = ?", info.Status)
	}

	// 创建时间范围搜索
	if info.StartDate != "" {
		db = db.Where("created_at >= ?", info.StartDate)
	}
	if info.EndDate != "" {
		db = db.Where("created_at <= ?", info.EndDate+" 23:59:59")
	}

	// 关键字搜索（同时搜索应用名称、应用ID、开发者名称、描述）
	if info.Keyword != "" {
		keyword := "%" + info.Keyword + "%"
		db = db.Where("app_name LIKE ? OR app_id LIKE ? OR description LIKE ?",
			keyword, keyword, keyword, keyword)
	}

	return db
}

// buildOrderConditions 构建排序条件
func (a *ApplicationService) buildOrderConditions(order string, desc bool) string {
	// 默认排序
	defaultOrder := "sort_order DESC"

	if order == "" {
		return defaultOrder
	}

	// 验证排序字段安全性
	allowedOrderFields := map[string]bool{
		"id":          true,
		"app_id":      true,
		"app_name":    true,
		"category_id": true,
		"status":      true,
		"created_at":  true,
		"updated_at":  true,
		"sort_order":  true,
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
	if order != "created_at" {
		orderStr += ", created_at DESC"
	}

	return orderStr
}
