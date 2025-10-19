package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/common/response"
	"ApkAdmin/model/project/request"
	"github.com/gin-gonic/gin"
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
	err = websiteConfigService.SetConfig(req.Scope, req.Config)
	if err != nil {
		global.GVA_LOG.Error("站点配置失败", zap.Error(err))
		response.FailWithMessage("站点配置失败，"+err.Error(), c)
		return
	}
	response.OkWithMessage("配置成功", c)
	return
}
