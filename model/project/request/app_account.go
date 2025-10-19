package request

import "ApkAdmin/constants"

// CreateAppAccountRequest 创建账号请求
type CreateAppAccountRequest struct {
	AppID         string                     `json:"app_id" binding:"required"`
	AccountDetail string                     `json:"account_detail" binding:"required"`
	ExtraInfo     interface{}                `json:"extra_info"`
	AccountStatus constants.AppAccountStatus `json:"account_status" binding:"omitempty,min=1,max=4"` // 1-4
}

// UpdateAppAccountRequest 更新账号请求
type UpdateAppAccountRequest struct {
	ID            uint                       `json:"id" binding:"required"`
	AppID         string                     `json:"app_id" binding:"required"`
	AccountDetail string                     `json:"account_detail"`
	ExtraInfo     interface{}                `json:"extra_info"`
	AccountStatus constants.AppAccountStatus `json:"account_status" binding:"omitempty,min=1,max=5"`
}

// GetAppAccountListRequest 获取账号列表请求
type GetAppAccountListRequest struct {
	PageInfo
	AppID         string `json:"app_id" form:"app_id"`
	CategoryID    uint   `json:"category_id" form:"category_id"`
	AccountNo     string `json:"account_no" form:"account_no"`
	AccountStatus int    `json:"account_status" form:"account_status"` // 0表示全部
}

// UpdateAccountStatusRequest 批量更新状态请求
type UpdateAccountStatusRequest struct {
	IDs    []uint                     `json:"ids" binding:"required,min=1"`
	Status constants.AppAccountStatus `json:"status" binding:"required,min=1,max=5"`
}

type ViewAppAccountOrderRequest struct {
	AccountId uint `json:"account_id" form:"account_id" binding:"required"`
}
