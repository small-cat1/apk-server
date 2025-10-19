package project

import (
	"ApkAdmin/middleware"
	"github.com/gin-gonic/gin"
)

type AppAccountRouter struct {
}

func (r *AppAccountRouter) InitAppAccountRouter(Router *gin.RouterGroup) {
	router := Router.Group("appAccount").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("appAccount")
	{
		router.POST("", appAccountApi.CreateAppAccount)                  // 创建应用账号
		router.PUT("", appAccountApi.UpdateAppAccount)                   // 更新应用账号
		router.PUT("batchUpdateStatus", appAccountApi.BatchUpdateStatus) // 批量更新应用账号状态
		router.DELETE("", appAccountApi.DeleteAppAccount)                // 删除应用账号

	}
	{
		routerWithoutRecord.GET("list", appAccountApi.ListAppAccount)        // 应用账号列表
		routerWithoutRecord.GET("detail", appAccountApi.GetAppAccountDetail) // 应用账号详情
		routerWithoutRecord.GET("order", appAccountApi.ViewAppAccountOrder)  // 查看售出账号关联订单

	}
}
