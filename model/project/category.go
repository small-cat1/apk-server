package project

import (
	"time"
)

const (
	// 顶级分类父ID
	TopLevelParentID = 0
)

// AppCategory 应用分类表结构体
type AppCategory struct {
	ID            uint      `json:"id" gorm:"primarykey;comment:分类ID"`
	ParentID      uint      `json:"parent_id" gorm:"default:0;index:idx_parent_id;comment:父分类ID（0表示顶级分类）"`
	CategoryCode  string    `json:"category_code" gorm:"size:50;uniqueIndex:uk_category_code;not null;comment:分类代码（唯一标识）"`
	CategoryName  string    `json:"category_name" gorm:"size:100;not null;comment:分类名称"`
	EmojiIcon     *string   `json:"emoji_icon" gorm:"size:100;comment:emoji图标"`
	Icon          *string   `json:"icon" gorm:"size:255;comment:分类图标"`
	Description   *string   `json:"description" gorm:"type:text;comment:分类描述"`
	SortOrder     int       `json:"sort_order" gorm:"default:0;index:idx_sort_order;comment:排序权重"`
	IsActive      *int      `json:"is_active" gorm:"default:1;comment:应用列表状态"`
	AccountStatus *int      `json:"account_status" gorm:"default:1;comment:账号列表状态"`
	TrendingTag   *int      `json:"trending_tag" gorm:"default:0;comment:是否是热门标签"` // 是否是热门标签
	IsBanner      *int      `json:"is_banner" gorm:"default:0;comment:是否是幻灯片"`     // 是否是幻灯片
	BannerUrl     *string   `json:"banner_url" gorm:"comment:幻灯片图片URL"`            // 幻灯片图片URL
	CreatedAt     time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"comment:更新时间"`
	// 关联字段（不存储在数据库中）
	Children []*AppCategory `json:"children,omitempty" gorm:"-"`                                // 子分类
	Parent   *AppCategory   `json:"parent,omitempty" gorm:"foreignKey:ParentID;references:ID;"` // 父分类
	// 新增：应用数量统计字段
	AppCount int64 `json:"app_count" gorm:"-"` // 当前分类下的应用数量
}

// TableName 指定表名
func (AppCategory) TableName() string {
	return "app_categories"
}

// IsTopLevel 判断是否为顶级分类
func (c *AppCategory) IsTopLevel() bool {
	return c.ParentID == TopLevelParentID
}

// ValidateCircularReference 验证是否存在循环引用
func ValidateCircularReference(categories []AppCategory, categoryID, newParentID uint) bool {
	if categoryID == newParentID {
		return false // 自己不能是自己的父级
	}

	// 查找新父级的所有父级路径
	var hasCircular func(uint) bool
	hasCircular = func(parentID uint) bool {
		if parentID == categoryID {
			return false // 发现循环引用
		}

		if parentID == TopLevelParentID {
			return true // 到达顶级，无循环
		}

		// 查找父级的父级
		for _, cat := range categories {
			if cat.ID == parentID {
				return hasCircular(cat.ParentID)
			}
		}

		return true
	}

	return hasCircular(newParentID)
}
