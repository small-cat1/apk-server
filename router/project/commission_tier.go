package project

import (
	"ApkAdmin/middleware"
	"github.com/gin-gonic/gin"
)

// CommissionTierRouter 分佣等级路由
type CommissionTierRouter struct {
}

func (r CommissionTierRouter) InitCommissionTierRouter(Router *gin.RouterGroup) {
	router := Router.Group("commissionTier").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("commissionTier")
	{
		router.POST("", commissionTierApi.CreateCommissionTier)
		router.PUT("", commissionTierApi.UpdateCommissionTier)
		router.DELETE("", commissionTierApi.DeleteCommissionTiers)
		router.PUT("updateStatus", commissionTierApi.UpdateCommissionTierStatus)
		router.PUT("updateSort", commissionTierApi.UpdateCommissionTierSort)
	}
	{
		routerWithoutRecord.GET("list", commissionTierApi.GetCommissionTierList)
		routerWithoutRecord.GET("detail", commissionTierApi.GetCommissionTier)
	}
}
