package project

import (
	"ApkAdmin/middleware"
	"github.com/gin-gonic/gin"
)

type PaymentAccountRouter struct {
}

func (r *PaymentAccountRouter) InitPaymentAccountRouter(Router *gin.RouterGroup) {
	router := Router.Group("paymentAccounts").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("paymentAccounts")
	{
		// 写操作路由 - 需要记录操作日志
		router.POST("", paymentAccountApi.CreatePaymentAccount)       // 创建支付账号
		router.PUT("", paymentAccountApi.UpdatePaymentAccount)        // 更新支付账号
		router.DELETE("", paymentAccountApi.DeletePaymentAccount)     // 删除支付账号，软删除
		router.PUT("weight", paymentAccountApi.UpdateAccountWeight)   // 更新账号权重
		router.PUT("reset-daily", paymentAccountApi.ResetDailyAmount) // 重置日交易金额
	}
	{
		// 读操作路由 - 不需要记录操作日志
		routerWithoutRecord.GET("", paymentAccountApi.GetPaymentAccountList)   // 获取支付账号列表
		routerWithoutRecord.GET("detail", paymentAccountApi.GetPaymentAccount) // 获取单个支付账号详情

	}
}
