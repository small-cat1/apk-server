package request

import (
	"ApkAdmin/constants"
	"ApkAdmin/model/common/request"
	"ApkAdmin/model/project"
	"encoding/json"
	"errors"
	"time"
)

// MembershipPlanListRequest 会员套餐列表请求
type MembershipPlanListRequest struct {
	request.PageInfo
	PlanCode     string   `json:"plan_code" form:"plan_code"`         // 套餐代码（模糊搜索）
	PlanName     string   `json:"plan_name" form:"plan_name"`         // 套餐名称（模糊搜索）
	PlanType     string   `json:"plan_type" form:"plan_type"`         // 套餐类型（精确匹配）
	Platform     string   `json:"platform" form:"platform"`           // 支持平台（JSON数组查询）
	CurrencyCode string   `json:"currency_code" form:"currency_code"` // 货币代码
	IsActive     *bool    `json:"is_active" form:"is_active"`         // 是否启用
	IsFeatured   *bool    `json:"is_featured" form:"is_featured"`     // 是否推荐
	MinPrice     *float64 `json:"min_price" form:"min_price"`         // 最低价格
	MaxPrice     *float64 `json:"max_price" form:"max_price"`         // 最高价格
	StartDate    string   `json:"start_date" form:"start_date"`       // 创建开始日期
	EndDate      string   `json:"end_date" form:"end_date"`           // 创建结束日期
	Keyword      string   `json:"keyword" form:"keyword"`             // 关键字搜索
}

// MembershipPlanCreateRequest 创建会员套餐请求结构体
type MembershipPlanCreateRequest struct {
	PlanCode             string             `json:"plan_code" validate:"required,min=2,max=50,alphanum_underscore" binding:"required"`
	PlanName             string             `json:"plan_name" validate:"required,min=1,max=100" binding:"required"`
	PlanType             constants.PlanType `json:"plan_type" validate:"required,oneof=monthly yearly lifetime" binding:"required"`
	Platform             json.RawMessage    `json:"platform" validate:"required,oneof=android ios harmony windows" binding:"required"`
	DurationDays         *int               `json:"duration_days" validate:"omitempty,min=1"`
	BasePrice            float64            `json:"base_price" validate:"required,gte=0" binding:"required"`
	CurrencyCode         string             `json:"currency_code" validate:"required,len=3,uppercase" binding:"required"`
	DiscountPercentage   float64            `json:"discount_percentage" validate:"gte=0,lte=100"`
	FinalPrice           float64            `json:"final_price" validate:"required,gte=0" binding:"required"`
	DownloadLimitDaily   *int               `json:"download_limit_daily" validate:"omitempty,min=0"`
	DownloadLimitMonthly *int               `json:"download_limit_monthly" validate:"omitempty,min=0"`
	IsActive             *bool              `json:"is_active"`
	IsFeatured           *bool              `json:"is_featured"`
	SortOrder            int                `json:"sort_order" validate:"gte=0"`
	Description          string             `json:"description" validate:"max=1000"`
}

// MembershipPlanUpdateRequest 更新会员套餐请求结构体
type MembershipPlanUpdateRequest struct {
	ID                   uint               `json:"id" validate:"required,gt=0" binding:"required"`
	PlanCode             string             `json:"plan_code" validate:"omitempty,min=2,max=50,alphanum_underscore"`
	PlanName             string             `json:"plan_name" validate:"omitempty,min=1,max=100"`
	PlanType             constants.PlanType `json:"plan_type" validate:"omitempty,oneof=monthly yearly lifetime"`
	Platform             json.RawMessage    `json:"platform" validate:"omitempty,oneof=android ios harmony windows"`
	DurationDays         *int               `json:"duration_days" validate:"omitempty,min=1"`
	BasePrice            float64            `json:"base_price" validate:"omitempty,gte=0"`
	CurrencyCode         string             `json:"currency_code" validate:"omitempty,len=3,uppercase"`
	DiscountPercentage   float64            `json:"discount_percentage" validate:"omitempty,gte=0,lte=100"`
	FinalPrice           float64            `json:"final_price" validate:"omitempty,gte=0"`
	DownloadLimitDaily   *int               `json:"download_limit_daily" validate:"omitempty,min=0"`
	DownloadLimitMonthly *int               `json:"download_limit_monthly" validate:"omitempty,min=0"`
	IsActive             *bool              `json:"is_active"`
	IsFeatured           *bool              `json:"is_featured"`
	SortOrder            *int               `json:"sort_order" validate:"omitempty,gte=0"`
	Description          *string            `json:"description" validate:"omitempty,max=1000"`
}

// Validate 验证创建请求的业务逻辑
func (req *MembershipPlanCreateRequest) Validate() error {
	// 终身套餐不应该有有效期
	if req.PlanType == "lifetime" && req.DurationDays != nil {
		return errors.New("终身套餐不应该设置有效天数")
	}

	// 非终身套餐必须有有效期
	if req.PlanType != "lifetime" && (req.DurationDays == nil || *req.DurationDays <= 0) {
		return errors.New("非终身套餐必须设置有效天数")
	}

	// 验证最终价格是否合理（应该小于等于基础价格）
	if req.FinalPrice > req.BasePrice {
		return errors.New("最终价格不能大于基础价格")
	}

	// 验证折扣和价格的一致性
	if req.DiscountPercentage > 0 {
		expectedPrice := req.BasePrice * (1 - req.DiscountPercentage/100)
		if abs(req.FinalPrice-expectedPrice) > 0.01 { // 允许0.01的误差
			return errors.New("最终价格与折扣计算不一致")
		}
	}

	// 验证下载限制的合理性
	if req.DownloadLimitDaily != nil && req.DownloadLimitMonthly != nil {
		if *req.DownloadLimitDaily*30 > *req.DownloadLimitMonthly {
			return errors.New("每日下载限制过高，超出月限制")
		}
	}

	return nil
}

// ValidateUpdateRequest 验证更新请求的业务逻辑
func (req *MembershipPlanUpdateRequest) Validate() error {
	// 如果更新了套餐类型，需要验证有效期设置
	if req.PlanType != "" {
		if req.PlanType == "lifetime" && req.DurationDays != nil {
			return errors.New("终身套餐不应该设置有效天数")
		}
		if req.PlanType != "lifetime" && req.DurationDays != nil && *req.DurationDays <= 0 {
			return errors.New("非终身套餐必须设置正确的有效天数")
		}
	}

	// 如果同时更新了基础价格和最终价格，验证合理性
	if req.BasePrice >= 0 && req.FinalPrice >= 0 {
		if req.FinalPrice > req.BasePrice {
			return errors.New("最终价格不能大于基础价格")
		}
	}

	// 验证折扣和价格的一致性（如果都提供了）
	if req.BasePrice >= 0 && req.FinalPrice >= 0 && req.DiscountPercentage != 0 {
		if req.DiscountPercentage > 0 {
			expectedPrice := req.BasePrice * (1 - req.DiscountPercentage/100)
			if abs(req.FinalPrice-expectedPrice) > 0.01 {
				return errors.New("最终价格与折扣计算不一致")
			}
		}
	}

	// 验证下载限制
	if req.DownloadLimitDaily != nil && req.DownloadLimitMonthly != nil {
		if *req.DownloadLimitDaily*30 > *req.DownloadLimitMonthly {
			return errors.New("每日下载限制过高，超出月限制")
		}
	}

	return nil
}

// ToMembershipPlan 将创建请求转换为MembershipPlan实体
func (req *MembershipPlanCreateRequest) ToMembershipPlan() *project.MembershipPlan {
	plan := &project.MembershipPlan{
		PlanCode:             req.PlanCode,
		PlanName:             req.PlanName,
		PlanType:             req.PlanType,
		Platform:             req.Platform,
		DurationDays:         req.DurationDays,
		BasePrice:            &req.BasePrice,
		CurrencyCode:         req.CurrencyCode,
		DiscountPercentage:   &req.DiscountPercentage,
		FinalPrice:           &req.FinalPrice,
		DownloadLimitDaily:   req.DownloadLimitDaily,
		DownloadLimitMonthly: req.DownloadLimitMonthly,
		SortOrder:            &req.SortOrder,
		Description:          stringPointer(req.Description),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	// 设置默认值
	if req.IsActive != nil {
		plan.IsActive = req.IsActive
	} else {
		active := true
		plan.IsActive = &active // 默认启用
	}

	if req.IsFeatured != nil {
		plan.IsFeatured = req.IsFeatured
	} else {
		featured := false
		plan.IsFeatured = &featured // 默认不推荐
	}

	// 如果是终身套餐，确保DurationDays为nil
	if req.PlanType == "lifetime" {
		plan.DurationDays = nil
	}

	return plan
}

// UpdateMembershipPlan 根据更新请求更新实体
func (req *MembershipPlanUpdateRequest) UpdateMembershipPlan(plan *project.MembershipPlan) {
	if req.PlanCode != "" {
		plan.PlanCode = req.PlanCode
	}
	if req.PlanName != "" {
		plan.PlanName = req.PlanName
	}
	if req.PlanType != "" {
		plan.PlanType = req.PlanType
		// 如果改为终身套餐，清空有效期
		if req.PlanType == "lifetime" {
			plan.DurationDays = nil
		}
	}
	if req.Platform != nil {
		plan.Platform = req.Platform
	}
	if req.DurationDays != nil {
		plan.DurationDays = req.DurationDays
	}
	if req.BasePrice >= 0 {
		plan.BasePrice = &req.BasePrice
	}
	if req.CurrencyCode != "" {
		plan.CurrencyCode = req.CurrencyCode
	}
	if req.DiscountPercentage >= 0 {
		plan.DiscountPercentage = &req.DiscountPercentage
	}
	if req.FinalPrice >= 0 {
		plan.FinalPrice = &req.FinalPrice
	}
	if req.DownloadLimitDaily != nil {
		plan.DownloadLimitDaily = req.DownloadLimitDaily
	}
	if req.DownloadLimitMonthly != nil {
		plan.DownloadLimitMonthly = req.DownloadLimitMonthly
	}
	if req.IsActive != nil {
		plan.IsActive = req.IsActive
	}
	if req.IsFeatured != nil {
		plan.IsFeatured = req.IsFeatured
	}
	if req.SortOrder != nil {
		plan.SortOrder = req.SortOrder
	}
	if req.Description != nil {
		plan.Description = req.Description
	}

	plan.UpdatedAt = time.Now()
}

// 切换状态请求结构体
type MembershipPlanToggleStatusRequest struct {
	ID       uint `json:"id" validate:"required,gt=0" binding:"required"`
	IsActive bool `json:"is_active"`
}

// 设置推荐请求结构体
type MembershipPlanSetFeaturedRequest struct {
	ID         uint `json:"id" validate:"required,gt=0" binding:"required"`
	IsFeatured bool `json:"is_featured"`
}

// 更新排序请求结构体
type MembershipPlanUpdateSortRequest struct {
	ID        uint `json:"id" validate:"required,gt=0" binding:"required"`
	SortOrder int  `json:"sort_order" validate:"gte=0"`
}

// 工具函数
func stringPointer(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
