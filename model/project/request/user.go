package request

import (
	"ApkAdmin/constants"
	"ApkAdmin/model/common/request"
	"ApkAdmin/utils"
	"errors"
)

// UserListRequest 用户列表请求参数
type UserListRequest struct {
	request.PageInfo
	// 分页信息
	Username      string                  `json:"username" form:"username"`             // 用户名（模糊搜索）
	Email         string                  `json:"email" form:"email"`                   // 邮箱（模糊搜索）
	Phone         string                  `json:"phone" form:"phone"`                   // 手机号（模糊搜索）
	AccountStatus constants.AccountStatus `json:"account_status" form:"account_status"` // 账户状态
	HasMembership string                  `json:"has_membership" form:"has_membership"` // 是否有会员
	Gender        string                  `json:"gender" form:"gender"`                 // 性别
	StartDate     string                  `json:"start_date" form:"start_date"`         // 注册开始日期
	EndDate       string                  `json:"end_date" form:"end_date"`             // 注册结束日期
	Keyword       string                  `json:"keyword" form:"keyword"`               // 关键字搜索
}

// UserRequest 用户详情请求参数
type UserRequest struct {
	ID uint `json:"id" form:"id" uri:"id"` // 用户ID
}

// CreateUserRequest 创建用户请求参数
type CreateUserRequest struct {
	Email         string `json:"email" binding:"required,email"`    // 邮箱
	Phone         string `json:"phone"`                             // 手机号
	Password      string `json:"password" binding:"required,min=6"` // 密码
	EmailVerified bool   `json:"email_verified"`                    // 邮箱验证状态
	PhoneVerified bool   `json:"phone_verified"`                    // 手机验证状态
}

// UpdateUserRequest 更新用户请求参数
type UpdateUserRequest struct {
	ID            uint   `json:"id" binding:"required"`                                                       // 用户ID
	Email         string `json:"email" binding:"required,email"`                                              // 邮箱
	Phone         string `json:"phone"`                                                                       // 手机号
	AccountStatus string `json:"account_status" binding:"oneof=active suspended pending_verification banned"` // 账户状态
	EmailVerified bool   `json:"email_verified"`                                                              // 邮箱验证状态
	PhoneVerified bool   `json:"phone_verified"`                                                              // 手机验证状态
}

// BatchDeleteUsersRequest 批量删除用户请求参数
type BatchDeleteUsersRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"` // 用户ID列表
}

// BatchUpdateUserStatusRequest 批量更新用户状态请求参数
type BatchUpdateUserStatusRequest struct {
	IDs           []uint `json:"ids" binding:"required,min=1"`                                    // 用户ID列表
	AccountStatus string `json:"account_status" binding:"required,oneof=active suspended banned"` // 账户状态
}

// ResetPasswordRequest 重置密码请求参数
type ResetPasswordRequest struct {
	ID uint `json:"id" binding:"required"` // 用户ID
}

// GetUserMembershipsRequest 获取用户会员记录请求参数
type GetUserMembershipsRequest struct {
	UserID uint `json:"user_id" form:"user_id" uri:"user_id" binding:"required"` // 用户ID
}

// GetUserOrdersRequest 获取用户订单记录请求参数
type GetUserOrdersRequest struct {
	UserID uint `json:"user_id" form:"user_id" uri:"user_id" binding:"required"` // 用户ID
}

type ChangeUserPasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

func (c ChangeUserPasswordRequest) Validate() error {
	if c.OldPassword == c.NewPassword {
		return errors.New("新旧密码不能一样！")
	}
	// 验证密码
	if err := utils.ValidatePassword(c.NewPassword); err != nil {
		return err
	}
	return nil
}
