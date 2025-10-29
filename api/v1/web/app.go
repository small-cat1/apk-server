package web

import (
	"ApkAdmin/global"
	"ApkAdmin/model/common/response"
	"ApkAdmin/model/project/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AppApi struct {
}

func (a AppApi) ListHotOrRecommendApp(c *gin.Context) {
	appList, err := AppService.GetHotOrRecommendApp()
	if err != nil {
		global.GVA_LOG.Error("获取热门或推荐应用失败!", zap.Error(err))
		response.FailWithMessage("获取热门或推荐应用失败，"+err.Error(), c)
		return
	}
	response.OkWithDetailed(appList, "获取成功", c)
}

func (a AppApi) GetFilterApps(c *gin.Context) {
	var req request.FilterAppRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		global.GVA_LOG.Error("获取分类列表应用失败!", zap.Error(err))
		response.FailWithMessage("获取分类列表应用失败", c)
		return
	}
	err = req.Validate()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := AppService.FilterAppsByCategory(req)
	if err != nil {
		global.GVA_LOG.Error("获取分类列表应用失败!", zap.Error(err))
		response.FailWithMessage("获取分类列表应用失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取分类列表应用成功", c)
}

// GetAccountAppsListByCategory 根据分类获取应用账号
func (a AppApi) GetAccountAppsListByCategory(c *gin.Context) {
	var req request.FilterAccountAppRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		global.GVA_LOG.Error("获取账号应用列表失败!", zap.Error(err))
		response.FailWithMessage("获取账号应用列表失败", c)
		return
	}
	err = req.Validate()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := AppService.GetAppByAccountCategory(req)
	if err != nil {
		global.GVA_LOG.Error("获取账号应用列表失败!", zap.Error(err))
		response.FailWithMessage("获取账号应用列表失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取账号应用列表成功", c)
}

func (a AppApi) SearchApps(c *gin.Context) {
	var req request.SearchAppRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		global.GVA_LOG.Error("搜索应用失败!", zap.Error(err))
		response.FailWithMessage("搜索应用失败", c)
		return
	}
	applicationList, total, err := AppService.SearchApps(req)
	if err != nil {
		global.GVA_LOG.Error("搜索应用失败!", zap.Error(err))
		response.FailWithMessage("搜索应用失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(map[string]interface{}{
		"list":  applicationList,
		"total": total,
	}, "获取成功", c)
}
