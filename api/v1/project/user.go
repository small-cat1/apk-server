package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/common/response"
	"ApkAdmin/model/project/request"
	"ApkAdmin/service/project"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserApi struct{}

// GetUserList 获取用户列表
func (u *UserApi) GetUserList(c *gin.Context) {
	var req request.UserListRequest

	// 绑定查询参数
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 获取排序参数
	order := c.DefaultQuery("order", "")
	desc := c.DefaultQuery("desc", "false") == "true"

	// 调用服务层
	if list, total, err := UserService.GetUserList(req, order, desc); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		}, "获取成功", c)
	}
}

// GetUser 获取用户详情
func (u *UserApi) GetUser(c *gin.Context) {
	var req request.UserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if user, err := UserService.GetUserDetail(project.WithID(req.ID)); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithData(gin.H{"user": user}, c)
	}
}

// DeleteUser 删除用户
func (u *UserApi) DeleteUser(c *gin.Context) {
	var req request.UserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := UserService.DeleteUser(req.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// BatchUpdateUserStatus 批量更新用户状态
func (u *UserApi) BatchUpdateUserStatus(c *gin.Context) {
	var req request.BatchUpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := UserService.BatchUpdateUserStatus(req.IDs, req.AccountStatus); err != nil {
		global.GVA_LOG.Error("批量更新状态失败!", zap.Error(err))
		response.FailWithMessage("批量更新状态失败", c)
	} else {
		response.OkWithMessage("批量更新状态成功", c)
	}
}

// ResetUserPassword 重置用户密码
func (u *UserApi) ResetUserPassword(c *gin.Context) {
	var req request.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if newPassword, err := UserService.ResetUserPassword(req.ID); err != nil {
		global.GVA_LOG.Error("重置密码失败!", zap.Error(err))
		response.FailWithMessage("重置密码失败", c)
	} else {
		response.OkWithData(gin.H{"newPassword": newPassword}, c)
	}
}

// GetUserMemberships 获取用户会员记录
func (u *UserApi) GetUserMemberships(c *gin.Context) {
	var req request.GetUserMembershipsRequest

	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if memberships, err := UserService.GetUserMemberships(req.UserID); err != nil {
		global.GVA_LOG.Error("获取会员记录失败!", zap.Error(err))
		response.FailWithMessage("获取会员记录失败", c)
	} else {
		response.OkWithData(gin.H{"memberships": memberships}, c)
	}
}

// GetUserOrders 获取用户订单记录
func (u *UserApi) GetUserOrders(c *gin.Context) {
	var req request.GetUserOrdersRequest

	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if orders, err := UserService.GetUserOrders(req.UserID); err != nil {
		global.GVA_LOG.Error("获取订单记录失败!", zap.Error(err))
		response.FailWithMessage("获取订单记录失败", c)
	} else {
		response.OkWithData(gin.H{"orders": orders}, c)
	}
}
