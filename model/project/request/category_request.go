package request

import (
	"ApkAdmin/model/common/request"
	"errors"
	"strings"
)

type CategoryPageInfo struct {
	request.PageInfo
	CategoryName  string `json:"categoryName" form:"categoryName"`
	ParentId      uint   `json:"parent_id" form:"parent_id"`
	AccountStatus uint   `json:"account_status " form:"account_status"`
}

type SelectCategoryRequest struct {
	IsActive      int  `json:"is_active" form:"is_active"` // 应用列表状态
	AccountStatus uint `json:"account_status " form:"account_status"`
}

// CreateCategoryRequest 创建分类请求结构体
type CreateCategoryRequest struct {
	ParentID      uint    `json:"parent_id" binding:"min=0"`                      // 父分类ID
	CategoryCode  string  `json:"category_code" binding:"required,min=2,max=50"`  // 分类代码
	CategoryName  string  `json:"category_name" binding:"required,min=2,max=100"` // 分类名称
	EmojiIcon     *string `json:"emoji_icon" binding:"omitempty,max=100"`         // emoji图标
	Icon          *string `json:"icon" binding:"omitempty,max=255"`               // 分类图标
	Description   *string `json:"description" binding:"omitempty"`                // 分类描述
	SortOrder     int     `json:"sort_order" binding:"min=0,max=9999"`            // 排序权重
	IsActive      int     `json:"is_active"`                                      // 应用列表状态
	AccountStatus int     `json:"account_status"`                                 // 账号列表状态
	TrendingTag   int     `json:"trending_tag"`                                   // 是否是热门标签
	IsBanner      int     `json:"is_banner"`                                      // 是否是幻灯片
	BannerUrl     *string `json:"banner_url" binding:"omitempty,max=255"`         // 幻灯片图片URL
}

// Validate ValidateCreate 验证创建请求
func (req *CreateCategoryRequest) Validate() error {
	if strings.TrimSpace(req.CategoryCode) == "" {
		return errors.New("分类代码不能为空")
	}
	if strings.TrimSpace(req.CategoryName) == "" {
		return errors.New("分类名称不能为空")
	}

	if len(req.CategoryCode) < 2 || len(req.CategoryCode) > 50 {
		return errors.New("分类代码长度必须在2-50个字符之间")
	}

	if len(req.CategoryName) < 2 || len(req.CategoryName) > 100 {
		return errors.New("分类名称长度必须在2-100个字符之间")
	}

	if req.Icon != nil && len(*req.Icon) > 255 {
		return errors.New("图标字段长度不能超过255个字符")
	}

	if req.SortOrder < 0 || req.SortOrder > 9999 {
		return errors.New("排序权重必须在0-9999之间")
	}

	// 验证幻灯片相关字段
	if req.IsBanner != 0 && req.IsBanner != 1 {
		return errors.New("是否是幻灯片字段只能为0或1")
	}

	// 如果启用幻灯片，必须提供幻灯片图片URL
	if req.IsBanner == 1 {
		if req.BannerUrl == nil || strings.TrimSpace(*req.BannerUrl) == "" {
			return errors.New("启用幻灯片时必须提供幻灯片图片URL")
		}
	}

	if req.BannerUrl != nil && len(*req.BannerUrl) > 255 {
		return errors.New("幻灯片图片URL长度不能超过255个字符")
	}

	return nil
}

// UpdateCategoryRequest 更新分类请求结构体
type UpdateCategoryRequest struct {
	ID            uint    `json:"id" binding:"required,min=1"`                    // 分类ID
	ParentID      uint    `json:"parent_id" binding:"min=0"`                      // 父分类ID
	CategoryCode  string  `json:"category_code" binding:"required,min=2,max=50"`  // 分类代码
	CategoryName  string  `json:"category_name" binding:"required,min=2,max=100"` // 分类名称
	EmojiIcon     *string `json:"emoji_icon" binding:"omitempty,max=100"`         // emoji图标
	Icon          *string `json:"icon" binding:"omitempty,max=255"`               // 分类图标
	Description   *string `json:"description" binding:"omitempty"`                // 分类描述
	SortOrder     int     `json:"sort_order" binding:"min=0,max=9999"`            // 排序权重
	IsActive      int     `json:"is_active"`                                      // 应用列表状态
	AccountStatus int     `json:"account_status"`                                 // 账号列表状态
	TrendingTag   int     `json:"trending_tag"`                                   // 是否是热门标签
	IsBanner      int     `json:"is_banner"`                                      // 是否是幻灯片
	BannerUrl     *string `json:"banner_url" binding:"omitempty,max=255"`         // 幻灯片图片URL
}

// Validate 验证更新请求
func (req *UpdateCategoryRequest) Validate() error {
	if req.ID == 0 {
		return errors.New("分类ID不能为空")
	}

	if req.ID == req.ParentID {
		return errors.New("不能将自己设置为父分类")
	}

	if strings.TrimSpace(req.CategoryCode) == "" {
		return errors.New("分类代码不能为空")
	}

	if strings.TrimSpace(req.CategoryName) == "" {
		return errors.New("分类名称不能为空")
	}

	if len(req.CategoryCode) < 2 || len(req.CategoryCode) > 50 {
		return errors.New("分类代码长度必须在2-50个字符之间")
	}

	if len(req.CategoryName) < 2 || len(req.CategoryName) > 100 {
		return errors.New("分类名称长度必须在2-100个字符之间")
	}

	if req.Icon != nil && len(*req.Icon) > 255 {
		return errors.New("图标字段长度不能超过255个字符")
	}

	if req.SortOrder < 0 || req.SortOrder > 9999 {
		return errors.New("排序权重必须在0-9999之间")
	}

	// 验证幻灯片相关字段
	if req.IsBanner != 0 && req.IsBanner != 1 {
		return errors.New("是否是幻灯片字段只能为0或1")
	}

	// 如果启用幻灯片，必须提供幻灯片图片URL
	if req.IsBanner == 1 {
		if req.BannerUrl == nil || strings.TrimSpace(*req.BannerUrl) == "" {
			return errors.New("启用幻灯片时必须提供幻灯片图片URL")
		}
	}

	if req.BannerUrl != nil && len(*req.BannerUrl) > 255 {
		return errors.New("幻灯片图片URL长度不能超过255个字符")
	}

	return nil
}
