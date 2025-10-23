package web

import (
	"ApkAdmin/global"
	"ApkAdmin/model/common/response"
	"ApkAdmin/model/project/request"
	"ApkAdmin/service/project"
	"ApkAdmin/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserApi struct {
}

func (u *UserApi) GetUserInfo(c *gin.Context) {
	uuid := utils.GetUserUuid(c)
	ReqUser, err := UserService.GetUserDetail(project.WithUuid(uuid))
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(gin.H{"userInfo": ReqUser}, "获取成功", c)
}

func (u *UserApi) ChangePassword(c *gin.Context) {
	var req request.ChangeUserPasswordRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("用户修改密码参数不正确！", zap.Error(err))
		response.OkWithMessage("修改密码参数不正确！", c)
		return
	}
	err = req.Validate()
	if err != nil {
		response.OkWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	err = UserService.ChangeUserPassword(userId, req)
	if err != nil {
		global.GVA_LOG.Error("用户密码修改失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("用户密码修改成功", c)
}

func (u *UserApi) Withdraw(c *gin.Context) {
	var req request.UserWithdrawRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("用户修改密码参数不正确！", zap.Error(err))
		response.OkWithMessage("修改密码参数不正确！", c)
		return
	}
	err = req.Validate()
	if err != nil {
		response.OkWithMessage(err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	err = UserService.ApplyWithdraw(userId, req)
	if err != nil {
		global.GVA_LOG.Error("用户密码修改失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithMessage("用户密码修改成功", c)
}

func (u *UserApi) GetWithdrawConfig(c *gin.Context) {
	config, err := systemConfigService.GetConfig("commission")
	if err != nil {
		global.GVA_LOG.Error("获取提现规则失败!", zap.Error(err))
		response.FailWithMessage("获取提现规则失败", c)
		return
	}
	response.OkWithData(config, c)
}
