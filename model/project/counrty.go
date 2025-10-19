package project

import (
	"ApkAdmin/model/common"
	"gorm.io/gorm"
	"strings"
	"time"
)

// CountryRegion 国家地区模型
type CountryRegion struct {
	ID                  uint             `json:"id" gorm:"primarykey;autoIncrement;comment:主键ID"`
	CountryCode         string           `json:"country_code" gorm:"column:country_code;type:varchar(3);uniqueIndex:uk_country_code;not null;comment:国家代码（ISO 3166-1）"`
	CountryName         string           `json:"country_name" gorm:"column:country_name;type:varchar(100);not null;comment:国家名称"`
	CountryNameEN       string           `json:"country_name_en" gorm:"column:country_name_en;type:varchar(100);not null;comment:英文国家名称"`
	Region              string           `json:"region" gorm:"column:region;type:varchar(50);comment:所属地区（如Asia, Europe）"`
	CurrencyCode        string           `json:"currency_code" gorm:"column:currency_code;type:varchar(3);comment:货币代码"`
	LanguageCodes       common.JSONSlice `json:"language_codes" gorm:"column:language_codes;type:json;comment:官方语言代码列表"`
	ContentRatingSystem string           `json:"content_rating_system" gorm:"column:content_rating_system;type:varchar(50);comment:内容分级系统（如ESRB, PEGI）"`
	IsSupported         *int             `json:"is_supported" gorm:"column:is_supported;type:tinyint(1);default:1;comment:是否支持分发"`
	CreatedAt           time.Time        `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt           time.Time        `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
}

func (CountryRegion) TableName() string {
	return "countries_regions"
}

// BeforeCreate GORM钩子 - 创建前
func (c *CountryRegion) BeforeCreate(tx *gorm.DB) error {
	// 确保代码为大写
	c.CountryCode = strings.ToUpper(c.CountryCode)
	if c.CurrencyCode != "" {
		c.CurrencyCode = strings.ToUpper(c.CurrencyCode)
	}
	return nil
}

// BeforeUpdate GORM钩子 - 更新前
func (c *CountryRegion) BeforeUpdate(tx *gorm.DB) error {
	// 确保代码为大写
	c.CountryCode = strings.ToUpper(c.CountryCode)
	if c.CurrencyCode != "" {
		c.CurrencyCode = strings.ToUpper(c.CurrencyCode)
	}
	return nil
}
