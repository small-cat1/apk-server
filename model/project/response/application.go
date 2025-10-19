package response

import "github.com/shopspring/decimal"

// AccountAppResp 账号应用响应response
type AccountAppResp struct {
	// "account_price": 8.99,
	//                "id": 6,
	//                "app_id": "5387a095-fdb1-4cd5-b85e-437ce035c59e",
	//                "app_name": "Youtube",
	//                "country_code": "US",
	//                "category_id": 6,
	//                "subcategory_id": null,
	//                "app_icon": "uploads/file/b63d76a9df0c598bb302b891ba5e7f33_20251015234326.jpg",
	//                "description": "下载适用于 Android 手机和平板电脑的 YouTube 官方应用。看世界之所看，享世界之所享 - 从最热门的音乐视频，到时下流行的游戏、时尚、美容、新闻和学习等类型的内容，全部尽揽眼底。您可以订阅喜爱的频道、创作自己的内容、与朋友分享精彩内容，还可以在任意设备上观看视频。\n",
	//                "is_hot": 1,
	//                "is_recommend": 1,
	//                "is_free": false,
	//                "rating": 4,
	//                "sort_order": 1,
	//                "status": "active",
	//                "created_at": "2025-10-15T23:44:53+08:00",
	//                "updated_at": "2025-10-15T23:45:05+08:00",
	//                "created_by": 1
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
