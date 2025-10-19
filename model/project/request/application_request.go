package request

import (
	"ApkAdmin/constants"
	"ApkAdmin/model/common/request"
	"ApkAdmin/model/project"
	"errors"
	"github.com/shopspring/decimal"
)

type ApplicationCreateRequest struct {
	AppName       string          `json:"app_name" binding:"required"`
	CountryCode   string          `json:"country_code" binding:"required"` // 国家代码
	CategoryID    *uint           `json:"category_id" binding:"required"`
	SubcategoryID *uint           `json:"subcategory_id"`
	AppIcon       *string         `json:"app_icon" binding:"required"`
	Description   *string         `json:"description"`
	IsHot         *int            `json:"is_hot" `
	IsRecommend   *int            `json:"is_recommend" `
	IsFree        bool            `json:"is_free" `
	SortOrder     int             `json:"sort_order" binding:"min=0,max=9999"` // 排序权重
	AccountPrice  decimal.Decimal `json:"account_price" binding:"required"`    //应用账号价格
}

func (r *ApplicationCreateRequest) Validate() error {
	if r.SortOrder < 0 || r.SortOrder > 9999 {
		return errors.New("排序权重必须在0-9999之间")
	}
	// 检查价格是否为空
	if r.AccountPrice.IsZero() {
		return errors.New("账号价格不能为0")
	}
	// 检查价格是否为负数
	if r.AccountPrice.IsNegative() {
		return errors.New("账号价格不能为负数")
	}
	// 检查价格范围
	minPrice := decimal.NewFromFloat(0.01)
	maxPrice := decimal.NewFromFloat(999999.99)
	if r.AccountPrice.LessThan(minPrice) {
		return errors.New("账号价格不能小于0.01")
	}
	if r.AccountPrice.GreaterThan(maxPrice) {
		return errors.New("账号价格不能大于999999.99")
	}
	// 检查小数位数
	if r.AccountPrice.Exponent() < -2 {
		return errors.New("账号价格最多保留两位小数")
	}
	// 自定义验证逻辑
	return nil
}

func (r *ApplicationCreateRequest) ToApplication() project.Application {
	return project.Application{
		AppName:      r.AppName,
		CountryCode:  r.CountryCode,
		AppIcon:      r.AppIcon,
		Description:  r.Description,
		AccountPrice: r.AccountPrice,
		IsHot:        r.IsHot,
		IsRecommend:  r.IsRecommend,
		IsFree:       &r.IsFree,
	}
}

type ApplicationUpdateRequest struct {
	ID uint `json:"id" binding:"required"`
	ApplicationCreateRequest
}

type ApplicationListRequest struct {
	request.PageInfo
	AppName       string `json:"app_name" form:"app_name"`
	AppID         string `json:"app_id" form:"app_id"`
	DeveloperName string `json:"developer_name" form:"developer_name"`
	CategoryID    int    `json:"category_id" form:"category_id"`
	Status        string `json:"status" form:"status"`
	StartDate     string `json:"start_date" form:"start_date"`
	EndDate       string `json:"end_date" form:"end_date"`
	Keyword       string `json:"keyword" form:"keyword"`
	OrderKey      string `json:"order_key" form:"order_key"`
	Desc          bool   `json:"desc" form:"desc"`
}

type ApplicationBatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required"`
}

type ApplicationBatchUpdateStatusRequest struct {
	IDs    []uint                      `json:"ids" binding:"required"`
	Status constants.ApplicationStatus `json:"status" binding:"required"`
}

type ApplicationCloneRequest struct {
	SourceID uint `json:"source_id" binding:"required"`
}

func (r *ApplicationCloneRequest) Validate() error {
	// 自定义验证逻辑
	return nil
}

type ApplicationUploadIconRequest struct {
	AppID uint64 `form:"app_id" binding:"required"`
}
