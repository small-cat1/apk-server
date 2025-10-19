package request

import (
	"ApkAdmin/model/common/request"
	"ApkAdmin/model/project"
	"fmt"
	"time"
)

// CountryCreateRequest 创建国家地区请求
type CountryCreateRequest struct {
	CountryCode         string   `json:"country_code" binding:"required" validate:"required"`
	CountryName         string   `json:"country_name" binding:"required,max=100" validate:"required,max=100"`
	CountryNameEN       string   `json:"country_name_en" binding:"required,max=100" validate:"required,max=100"`
	Region              string   `json:"region" binding:"max=50" validate:"max=50"`
	CurrencyCode        string   `json:"currency_code" binding:"omitempty,len=3" validate:"omitempty,len=3,alpha,uppercase"`
	LanguageCodes       []string `json:"language_codes"`
	ContentRatingSystem string   `json:"content_rating_system" binding:"max=50" validate:"max=50"`
	IsSupported         int      `json:"is_supported" binding:"oneof=0 1" validate:"oneof=0 1"`
}

// CountryUpdateRequest 更新国家地区请求
type CountryUpdateRequest struct {
	ID                  uint     `json:"id" binding:"required"`
	CountryCode         string   `json:"country_code" binding:"required" validate:"required"`
	CountryName         string   `json:"country_name" binding:"required,max=100" validate:"required,max=100"`
	CountryNameEN       string   `json:"country_name_en" binding:"required,max=100" validate:"required,max=100"`
	Region              string   `json:"region" binding:"max=50" validate:"max=50"`
	CurrencyCode        string   `json:"currency_code" binding:"omitempty,len=3" validate:"omitempty,len=3,alpha,uppercase"`
	LanguageCodes       []string `json:"language_codes"`
	ContentRatingSystem string   `json:"content_rating_system" binding:"max=50" validate:"max=50"`
	IsSupported         int      `json:"is_supported" binding:"oneof=0 1" validate:"oneof=0 1"`
}

// CountryListRequest 国家地区列表请求
type CountryListRequest struct {
	request.PageInfo
	Region string `form:"region"`
}

// CountryResponse 国家地区响应
type CountryResponse struct {
	ID                  uint      `json:"id"`
	CountryCode         string    `json:"country_code"`
	CountryName         string    `json:"country_name"`
	CountryNameEN       string    `json:"country_name_en"`
	Region              string    `json:"region"`
	CurrencyCode        string    `json:"currency_code"`
	LanguageCodes       []string  `json:"language_codes"`
	ContentRatingSystem string    `json:"content_rating_system"`
	IsSupported         int       `json:"is_supported"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

// 验证函数
func (r *CountryCreateRequest) Validate() error {
	// 自定义验证逻辑
	if r.CountryCode == "" {
		return fmt.Errorf("国家代码不能为空")
	}
	if len(r.CountryCode) < 2 || len(r.CountryCode) > 4 {
		return fmt.Errorf("国家代码必须为2-4位字符")
	}
	if r.CountryName == "" {
		return fmt.Errorf("国家名称不能为空")
	}
	if r.CountryNameEN == "" {
		return fmt.Errorf("英文国家名称不能为空")
	}
	if r.CurrencyCode != "" && len(r.CurrencyCode) != 3 {
		return fmt.Errorf("货币代码必须为3位字符")
	}
	return nil
}

func (r *CountryUpdateRequest) Validate() error {
	if r.ID == 0 {
		return fmt.Errorf("ID不能为空")
	}
	if r.CountryCode == "" {
		return fmt.Errorf("国家代码不能为空")
	}
	if len(r.CountryCode) < 2 || len(r.CountryCode) > 4 {
		return fmt.Errorf("国家代码必须为2-4位字符")
	}
	if r.CountryName == "" {
		return fmt.Errorf("国家名称不能为空")
	}
	if r.CountryNameEN == "" {
		return fmt.Errorf("英文国家名称不能为空")
	}
	if r.CurrencyCode != "" && len(r.CurrencyCode) != 3 {
		return fmt.Errorf("货币代码必须为3位字符")
	}
	return nil
}

// ToCountryRegion 转换为数据库模型
func (r *CountryCreateRequest) ToCountryRegion() *project.CountryRegion {
	return &project.CountryRegion{
		CountryCode:         r.CountryCode,
		CountryName:         r.CountryName,
		CountryNameEN:       r.CountryNameEN,
		Region:              r.Region,
		CurrencyCode:        r.CurrencyCode,
		LanguageCodes:       r.LanguageCodes,
		ContentRatingSystem: r.ContentRatingSystem,
		IsSupported:         &r.IsSupported,
	}
}
