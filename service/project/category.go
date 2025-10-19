package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	"ApkAdmin/model/project/request"
	"ApkAdmin/model/project/response"
	"errors"
	"gorm.io/gorm"
)

type CategoryService struct{}

func (a *CategoryService) GetCategory(id uint) (category project.AppCategory, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&category).Error
	return
}

func (a *CategoryService) CreateCategory(req request.CreateCategoryRequest) (err error) {
	e := project.AppCategory{
		ParentID:      req.ParentID,
		CategoryCode:  req.CategoryCode,
		CategoryName:  req.CategoryName,
		EmojiIcon:     req.EmojiIcon,
		Icon:          req.Icon,
		Description:   req.Description,
		SortOrder:     req.SortOrder,
		IsActive:      &req.IsActive,
		AccountStatus: &req.AccountStatus,
		TrendingTag:   &req.TrendingTag,
		IsBanner:      &req.IsBanner,
		BannerUrl:     req.BannerUrl,
	}
	err = global.GVA_DB.Create(&e).Error
	return err
}

func (a *CategoryService) UpdateCategory(req *request.UpdateCategoryRequest) (err error) {
	var category project.AppCategory
	err = global.GVA_DB.Where("category_code = ? and id != ?", req.CategoryCode, req.ID).First(&category).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if category.ID > 0 {
		return errors.New("存在相同的分类编码，请重新配置")
	}
	e := project.AppCategory{
		ParentID:      req.ParentID,
		CategoryCode:  req.CategoryCode,
		CategoryName:  req.CategoryName,
		EmojiIcon:     req.EmojiIcon,
		Icon:          req.Icon,
		Description:   req.Description,
		SortOrder:     req.SortOrder,
		IsActive:      &req.IsActive,
		AccountStatus: &req.AccountStatus,
		TrendingTag:   &req.TrendingTag,
		IsBanner:      &req.IsBanner,
		BannerUrl:     req.BannerUrl,
	}
	err = global.GVA_DB.Omit("created_at").Where("id = ?", req.ID).Updates(e).Error
	return err
}

func (a *CategoryService) DeleteCategory(cid int) (err error) {
	var total int64
	err = global.GVA_DB.Model(&project.Application{}).Where("category_id = ?", cid).Count(&total).Error
	if err != nil {
		return err
	}
	if total > 0 {
		return errors.New("该分类下存在应用，请先解除使用")
	}
	return global.GVA_DB.Model(&project.AppCategory{}).Where("id = ?", cid).Delete(&project.AppCategory{}).Error
}

func (a *CategoryService) FindSelectCategory(info request.SelectCategoryRequest) (list interface{}, err error) {
	var categoryLists []response.SelectCategoryResponse
	OrderStr := "sort_order desc"
	db := global.GVA_DB.Model(&project.AppCategory{}).Select("id,parent_id,category_code,category_name,sort_order,created_at")
	if info.IsActive == 1 {
		db = db.Where("is_active = ?", info.IsActive)
	}
	if info.AccountStatus == 1 {
		db = db.Where("account_status = ?", info.AccountStatus)
	}
	err = db.Order(OrderStr).Scan(&categoryLists).Error
	return categoryLists, err
}

func (a *CategoryService) GetCategoryList(info request.CategoryPageInfo, order string, desc bool) (list interface{}, total int64, err error) {
	var categoryLists []project.AppCategory
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&project.AppCategory{})

	if info.CategoryName != "" {
		db = db.Where("category_name LIKE ?", "%"+info.CategoryName+"%")
	}
	if info.ParentId > 0 {
		db = db.Where("parent_id = ?", info.ParentId)
	}
	err = db.Count(&total).Error
	if err != nil {
		return categoryLists, total, err
	}
	db = db.Select(`app_categories.*,
			(SELECT COUNT(*) 
			 FROM applications 
			 WHERE applications.category_id = app_categories.id 
			 AND applications.status = 'active') as app_count`).Limit(limit).Offset(offset)
	OrderStr := "sort_order desc"
	if order != "" {
		OrderStr = order
		if desc {
			OrderStr = order + " desc"
		}
	}
	err = db.Order(OrderStr).Scan(&categoryLists).Error
	if info.ParentId > 0 {
		return a.getChildrenList(categoryLists, info.ParentId), total, err
	} else {
		return a.getChildrenList(categoryLists, 0), total, err
	}
}

// getChildrenList 子类
func (a *CategoryService) getChildrenList(categories []project.AppCategory, parentID uint) []*project.AppCategory {
	var tree []*project.AppCategory
	for _, category := range categories {
		if category.ParentID == parentID {
			category.Children = a.getChildrenList(categories, category.ID)
			tree = append(tree, &category)
		}
	}
	return tree
}

// GetTrendingCategory 获取热门分类
func (a *CategoryService) GetTrendingCategory() (list interface{}, err error) {
	var categoryLists []response.CategoryResponse
	err = global.GVA_DB.Model(&project.AppCategory{}).
		Select(`app_categories.*,
			(SELECT COUNT(*) 
			 FROM applications 
			 WHERE applications.category_id = app_categories.id 
			 AND applications.status = 'active') as app_count`).
		Where("trending_tag = ?", 1).
		Order("sort_order desc").
		Scan(&categoryLists).Error
	if err != nil {
		return nil, err
	}
	return categoryLists, err
}

// GetAllCategory 获取所有分类
func (a *CategoryService) GetAllCategory() (list interface{}, err error) {
	var categoryLists []project.AppCategory
	err = global.GVA_DB.Model(&project.AppCategory{}).
		Where("parent_id = ? and is_active = ?", 0, 1).
		Order("sort_order desc").
		Scan(&categoryLists).Error
	if err != nil {
		return nil, err
	}
	return categoryLists, err
}

func (a *CategoryService) GetAccountListAppCategory() (list interface{}, err error) {
	var categoryLists []project.AppCategory
	err = global.GVA_DB.Model(&project.AppCategory{}).
		Where("parent_id = ? and account_status = ?", 0, 1).
		Order("sort_order desc").
		Scan(&categoryLists).Error
	if err != nil {
		return nil, err
	}
	return categoryLists, err
}
