package project

import (
	"ApkAdmin/constants"
	"encoding/json"
	"github.com/shopspring/decimal"
	"time"
)

// Application 应用基础信息表
type Application struct {
	ID                uint64                      `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	AppID             string                      `json:"app_id" gorm:"uniqueIndex:uk_app_id;type:varchar(100);not null;comment:应用唯一标识符"`
	AppName           string                      `json:"app_name" gorm:"type:varchar(200);not null;comment:应用名称"`
	CountryCode       string                      `json:"country_code" gorm:"size:3;comment:国家代码(NULL表示通用)"`
	CategoryID        *uint                       `json:"category_id" gorm:"index:idx_category_id;comment:应用分类ID"`
	SubcategoryID     *uint                       `json:"subcategory_id" gorm:"comment:子分类ID"`
	AppIcon           *string                     `json:"app_icon" gorm:"type:varchar(500);comment:应用图标URL"`
	Description       *string                     `json:"description" gorm:"type:text;comment:应用描述"`
	IsHot             *int                        `json:"is_hot" gorm:"default:0;comment:是否热门"`
	IsRecommend       *int                        `json:"is_recommend" gorm:"default:0;comment:是否推荐"`
	IsFree            *bool                       `json:"is_free" gorm:"default:0;comment:是否免费应用,0不是，1是"`
	Rating            *float64                    `json:"rating" gorm:"default:4;comment:评分"`
	DownloadCount     uint                        `json:"download_count" gorm:"default:0;comment:下载次数"`
	SalesCount        int64                       `json:"sales_count" gorm:"default:0;comment:总售卖次数"`
	ApkSalesCount     int64                       `json:"apk_sales_count" gorm:"default:0;comment:apk套餐售卖次数"`
	AccountSalesCount int64                       `json:"account_sales_count" gorm:"default:0;comment:账号售卖次数"`
	SortOrder         int                         `json:"sort_order" gorm:"default:0;index:idx_sort_order;comment:排序权重"`
	AccountPrice      decimal.Decimal             `json:"account_price" gorm:"default:4.00;comment:账号价格"`
	Status            constants.ApplicationStatus `json:"status" gorm:"type:enum('active','suspended','deleted');default:active;comment:应用状态"`
	CreatedAt         time.Time                   `json:"created_at" gorm:"autoCreateTime;comment:创建时间"`
	UpdatedAt         time.Time                   `json:"updated_at" gorm:"autoUpdateTime;comment:更新时间"`
	CreatedBy         int64                       `json:"created_by" gorm:"not null;comment:创建人ID"`
	Packages          []AppPackage                `json:"packages,omitempty" gorm:"foreignKey:AppID;references:AppID"`
	Accounts          []AppAccount                `json:"accounts,omitempty" gorm:"foreignKey:AppID;references:AppID"`
	// 主分类关联
	Category *AppCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID;references:ID"`
	// 子分类关联
	Subcategory *AppCategory `json:"subcategory,omitempty" gorm:"foreignKey:SubcategoryID;references:ID"`
}

// TableName 指定表名
func (Application) TableName() string {
	return "applications"
}

func (a Application) MarshalJSON() ([]byte, error) {
	type Alias Application
	return json.Marshal(&struct {
		AccountPrice float64 `json:"account_price"`
		*Alias
	}{
		AccountPrice: a.AccountPrice.InexactFloat64(),
		Alias:        (*Alias)(&a),
	})
}
