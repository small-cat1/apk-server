package response

import "time"

// AppAccountResponse 应用账号响应
type AppAccountResponse struct {
	ID                uint      `json:"id"`
	AppID             string    `json:"app_id"`
	AppName           string    `json:"app_name"`
	CategoryID        uint      `json:"category_id"`
	CategoryName      string    `json:"category_name"`
	AccountNo         string    `json:"account_no"`
	AccountDetail     string    `json:"account_detail,omitempty"` // 根据权限决定是否返回
	ExtraInfo         string    `json:"extra_info"`
	AccountStatus     int       `json:"account_status"`
	AccountStatusText string    `json:"account_status_text"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedBy         uint      `json:"created_by"`
	CreatorName       string    `json:"creator_name"`
}

// AccountStatusOption 状态选项
type AccountStatusOption struct {
	Value int    `json:"value"`
	Label string `json:"label"`
	Type  string `json:"type"` // Element Plus tag type
}

type FailDetail struct {
	AccountID uint   `json:"account_id"`
	Reason    string `json:"reason"`
}
type BatchUpdateStatusResponse struct {
	SuccessCount int          `json:"success_count"`
	FailCount    int          `json:"fail_count"`
	FailDetails  []FailDetail `json:"fail_details"`
}
