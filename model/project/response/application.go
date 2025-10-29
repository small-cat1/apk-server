package response

import "github.com/shopspring/decimal"

// AccountAppResp 账号应用响应response
type AccountAppResp struct {
	ID           uint64          `json:"id"`       // 账号ID
	AppID        string          `json:"app_id"`   // 应用ID
	AppName      string          `json:"app_name"` //应用名称
	AppIcon      *string         `json:"app_icon" //应用缩略图`
	CategoryID   *uint           `json:"category_id"`   //分类ID
	CategoryName string          `json:"category_name"` //分类名称
	AccountPrice decimal.Decimal `json:"account_price"` //单价
	Rating       *float64        `json:"rating" `       //评分
	Stock        int             `json:"stock"`         // 库存数量
	SalesCount   int64           `json:"sales_count" gorm:"default:0;comment:总售卖次数"`
}
