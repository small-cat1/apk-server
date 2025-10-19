package response

import "time"

// CategoryResponse 分类响应结构体
type CategoryResponse struct {
	ID           uint      `json:"id"`
	ParentID     uint      `json:"parent_id"`
	CategoryCode string    `json:"category_code"`
	CategoryName string    `json:"category_name"`
	EmojiIcon    *string   `json:"emoji_icon"`
	Icon         *string   `json:"icon"`
	Description  *string   `json:"description"`
	SortOrder    int       `json:"sort_order"`
	AppCount     int64     `json:"app_count"` // 应用数量
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type SelectCategoryResponse struct {
	ID           uint      `json:"id"`
	ParentID     uint      `json:"parent_id"`
	CategoryCode string    `json:"category_code"`
	CategoryName string    `json:"category_name"`
	CreatedAt    time.Time `json:"created_at"`
}
