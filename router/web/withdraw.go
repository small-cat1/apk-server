package web

import "github.com/gin-gonic/gin"

type WithdrawRouter struct {
}

func (r WithdrawRouter) InitWithdrawRouter(Router *gin.RouterGroup) {
	router := Router.Group("withdraw")
	router.POST("", withdrawApi.Withdraw)                // 用户提现
	router.GET("/config", withdrawApi.GetWithdrawConfig) // 获取用户提现配置
	router.GET("records", withdrawApi.GetRecords)        // 获取用户提现记录
}
