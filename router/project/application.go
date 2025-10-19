package project

import (
	"ApkAdmin/middleware"
	"github.com/gin-gonic/gin"
)

type ApplicationRouter struct {
}

// InitApplicationBasicRouter 应用基础管理路由
func (r *ApplicationRouter) InitApplicationBasicRouter(Router *gin.RouterGroup) {
	router := Router.Group("application").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("application")
	{
		// 需要记录操作日志的接口
		router.POST("create", applicationApi.CreateApplication)                        // 创建应用
		router.PUT("update", applicationApi.UpdateApplication)                         // 更新应用
		router.DELETE("delete", applicationApi.DeleteApplication)                      // 删除应用
		router.DELETE("batch-delete", applicationApi.BatchDeleteApplications)          // 批量删除应用
		router.PUT("batch-update-status", applicationApi.BatchUpdateApplicationStatus) // 批量更新应用状态
		router.POST("clone", applicationApi.CloneApplication)                          // 克隆应用
		router.POST("upload-icon", applicationApi.UploadApplicationIcon)               // 上传应用图标
	}
	{
		// 不需要记录操作日志的接口
		routerWithoutRecord.GET("list", applicationApi.GetApplicationList) // 应用列表
		routerWithoutRecord.GET("detail", applicationApi.GetApplication)   // 获取应用详情
	}
}
