package project

import (
	"ApkAdmin/middleware"
	"github.com/gin-gonic/gin"
)

type MembershipOrderRouter struct {
}

// InitMembershipOrderBasicRouter 会员订单管理路由
func (r *MembershipOrderRouter) InitMembershipOrderBasicRouter(Router *gin.RouterGroup) {
	router := Router.Group("membershipOrder").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("membershipOrder")
	{
		// 需要记录操作日志的接口
		router.PUT("updateMembershipOrderRemark", membershipOrderApi.UpdateMembershipOrderRemark) // 更新会员订单备注/标记
		router.PUT("cancelMembershipOrder", membershipOrderApi.CancelMembershipOrder)             // 取消会员订单
		router.PUT("batchCancelMembershipOrders", membershipOrderApi.BatchCancelMembershipOrders) // 批量取消会员订单
		router.PUT("refundMembershipOrder", membershipOrderApi.RefundMembershipOrder)             // 申请退款会员订单
		router.PUT("confirmPayment", membershipOrderApi.ConfirmPayment)                           // 手动确认支付
		router.POST("handlePaymentCallback", membershipOrderApi.HandlePaymentCallback)            // 处理支付回调
		router.POST("validateOrder", membershipOrderApi.ValidateOrder)                            // 验证订单有效性
		router.POST("manualProcessOrder", membershipOrderApi.ManualProcessOrder)                  // 手动处理异常订单
		router.POST("syncPaymentStatus", membershipOrderApi.SyncPaymentStatus)                    // 查询第三方支付状态
		router.POST("sendOrderNotification", membershipOrderApi.SendOrderNotification)            // 发送订单通知
	}
	{
		// 不需要记录操作日志的接口
		routerWithoutRecord.GET("getMembershipOrderList", membershipOrderApi.GetMembershipOrderList)            // 分页获取会员订单列表
		routerWithoutRecord.GET("findMembershipOrder", membershipOrderApi.GetMembershipOrder)                   // 用id查询会员订单
		routerWithoutRecord.GET("findMembershipOrderByOrderNo", membershipOrderApi.GetMembershipOrderByOrderNo) // 根据订单号查询会员订单
		routerWithoutRecord.GET("getOrderStats", membershipOrderApi.GetOrderStats)                              // 获取订单统计信息
		routerWithoutRecord.GET("getUserOrderHistory", membershipOrderApi.GetUserOrderHistory)                  // 获取用户订单历史
		routerWithoutRecord.GET("exportOrders", membershipOrderApi.ExportOrders)                                // 导出订单数据
		routerWithoutRecord.GET("getPaymentMethods", membershipOrderApi.GetPaymentMethods)                      // 获取支付方式列表
		routerWithoutRecord.GET("getOrderLogs", membershipOrderApi.GetOrderLogs)                                // 获取订单操作日志
		routerWithoutRecord.GET("getRefundDetail", membershipOrderApi.GetRefundDetail)                          // 获取退款详情
		routerWithoutRecord.GET("getOrderReceipt", membershipOrderApi.GetOrderReceipt)                          // 获取订单收据
	}
}
