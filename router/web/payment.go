package web

import "github.com/gin-gonic/gin"

type PaymentRouter struct {
}

func (r PaymentRouter) InitPaymentRouter(Router *gin.RouterGroup) {
	Router.GET("paymentMethods", paymentApi.GetPaymentProviders) //获取支付服务商
}
