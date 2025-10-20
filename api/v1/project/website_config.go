package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/common/response"
	"ApkAdmin/model/project/request"
	"ApkAdmin/utils"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
)

type WebsiteConfigApi struct {
}

func (w WebsiteConfigApi) GetSystemConfig(c *gin.Context) {
	var req request.GetSystemConfigRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = req.Validate()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	config, err := websiteConfigService.GetConfig(req.Scope)
	response.OkWithData(map[string]interface{}{
		"config": config,
	}, c)
	return
}

func (w WebsiteConfigApi) GetConfigByKey(c *gin.Context) {
	var req request.GetConfigByKeyRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = req.Validate()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	config, err := websiteConfigService.GetConfigByKey(req.Scope, req.Key)
	response.OkWithData(config, c)
}

func (w WebsiteConfigApi) SetSystemConfig(c *gin.Context) {
	var req request.SetSystemConfigRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = req.Validate()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.Scope == "commission" {
		// 获取当前用户
		uuid := utils.GetUserUuid(c)
		// 获取用户谷歌验证器密钥
		user, err := sysUserService.GetUserInfo(uuid)
		if err != nil {
			response.FailWithMessage("获取用户信息失败", c)
			return
		}
		if !user.GoogleAuthStatus {
			response.FailWithMessage("未绑定谷歌验证器,无法操作", c)
			return
		}
		// 验证谷歌验证码
		valid := totp.Validate(req.GoogleCode, user.GoogleAuthKey)
		if !valid {
			response.FailWithMessage("谷歌验证码不正确！", c)
			return
		}
	}
	err = websiteConfigService.SetConfig(req.Scope, req.Config)
	if err != nil {
		global.GVA_LOG.Error("站点配置失败", zap.Error(err))
		response.FailWithMessage("站点配置失败，"+err.Error(), c)
		return
	}
	response.OkWithMessage("配置成功", c)
	return
}
