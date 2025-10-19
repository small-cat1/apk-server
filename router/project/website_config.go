package project

import (
	"ApkAdmin/middleware"
	"github.com/gin-gonic/gin"
)

type WebsiteConfigRouter struct {
}

func (r *WebsiteConfigRouter) InitWebsiteConfigRouter(Router *gin.RouterGroup) {
	websiteRouter := Router.Group("website").Use(middleware.OperationRecord())
	websiteRouterWithoutRecord := Router.Group("website")
	{
		websiteRouter.POST("setSystemConfig", websiteConfigApi.SetSystemConfig)
	}
	{
		websiteRouterWithoutRecord.GET("getSystemConfig", websiteConfigApi.GetSystemConfig)
		websiteRouterWithoutRecord.GET("getConfigByKey", websiteConfigApi.GetConfigByKey)
	}
}
