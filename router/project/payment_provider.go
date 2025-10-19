package project

import (
	"ApkAdmin/middleware"
	"github.com/gin-gonic/gin"
)

type PaymentProviderRouter struct {
}

func (r *PaymentProviderRouter) InitPaymentProviderRouter(Router *gin.RouterGroup) {
	router := Router.Group("paymentProvider").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("paymentProvider")
	{
		// 写操作路由 - 需要记录操作日志
		router.POST("createPaymentProvider", paymentProviderApi.CreatePaymentProvider)    // 创建支付服务商
		router.PUT("updatePaymentProvider", paymentProviderApi.UpdatePaymentProvider)     // 更新支付服务商
		router.DELETE("deletePaymentProvider", paymentProviderApi.DeletePaymentProvider)  // 删除支付服务商
		router.POST("batch-status", paymentProviderApi.BatchUpdatePaymentProviderStatus)  // 批量更新支付服务商状态
		router.PUT("updateProviderSortOrder", paymentProviderApi.UpdateProviderSortOrder) // 更新服务商排序
	}
	{
		// 读操作路由 - 不需要记录操作日志
		routerWithoutRecord.GET("getPaymentProviderList", paymentProviderApi.GetPaymentProviderList) // 获取支付服务商列表
		routerWithoutRecord.GET("getPaymentProvider", paymentProviderApi.GetPaymentProvider)         // 获取单个支付服务商详情

	}
}
