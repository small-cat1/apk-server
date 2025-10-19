package web

import "github.com/gin-gonic/gin"

type OrderRouter struct {
}

func (r *OrderRouter) InitOrderRouter(Router *gin.RouterGroup) {
	router := Router.Group("order")
	router.POST("/account", orderApi.StoreAccountOrder)           // 账号商品下单
	router.POST("/membership", orderApi.StoreMembershipPlanOrder) //会员套餐下单
}
