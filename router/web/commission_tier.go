package web

import "github.com/gin-gonic/gin"

type CommissionTierRouter struct {
}

func (r CommissionTierRouter) InitCommissionTier(Router *gin.RouterGroup) {
	Router.GET("commissionTier", commissionTierApi.GetCommissionTier)
	return
}
