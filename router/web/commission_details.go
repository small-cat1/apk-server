package web

import (
	"github.com/gin-gonic/gin"
)

type CommissionDetailRouter struct{}

// InitCommissionDetailRouter 初始化分佣明细路由
func (r *CommissionDetailRouter) InitCommissionDetailRouter(Router *gin.RouterGroup) {
	commissionDetailRouterWithoutRecord := Router.Group("commissionDetail")
	{
		commissionDetailRouterWithoutRecord.GET("getCommissionDetailList", commissionDetailApi.GetCommissionDetailList) // 分页获取分佣明细列表（不记录操作日志）
	}
}
