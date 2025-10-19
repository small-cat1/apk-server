package project

import (
	"ApkAdmin/middleware"
	"github.com/gin-gonic/gin"
)

type AppPackageRouter struct {
}

func (r *AppPackageRouter) InitAppPackageRouter(Router *gin.RouterGroup) {
	router := Router.Group("appPackage").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("appPackage")
	{
		router.POST("create", appPackageApi.AddAppPackage)                           // 添加应用安装包
		router.PUT("update", appPackageApi.UpdateAppPackage)                         // 编辑应用安装包
		router.DELETE("delete", appPackageApi.DeleteAppPackage)                      // 删除应用安装包
		router.PUT("batch-update-status", appPackageApi.BatchUpdateAppPackageStatus) // 批量更新应用安装包状态
	}
	{
		routerWithoutRecord.GET("list", appPackageApi.GetAppPackageList) // 应用安装包列表
		routerWithoutRecord.GET("detail", appPackageApi.FirstAppPackage) // 获取单一应用安装包信息
	}
}
