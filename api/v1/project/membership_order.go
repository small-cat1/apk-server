package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/common/request"
	"ApkAdmin/model/common/response"
	"ApkAdmin/model/project"
	projectReq "ApkAdmin/model/project/request"
	"ApkAdmin/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MembershipOrderApi struct{}

// GetMembershipOrderList 分页获取会员订单列表
// @Tags MembershipOrder
// @Summary 分页获取会员订单列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取会员订单列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /membershipOrder/getMembershipOrderList [get]
func (m *MembershipOrderApi) GetMembershipOrderList(c *gin.Context) {
	var pageInfo projectReq.MembershipOrderSearchRequest
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := membershipOrderService.GetMembershipOrderInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetMembershipOrder 用id查询会员订单
// @Tags MembershipOrder
// @Summary 用id查询会员订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query project.MembershipOrder true "用id查询会员订单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /membershipOrder/findMembershipOrder [get]
func (m *MembershipOrderApi) GetMembershipOrder(c *gin.Context) {
	var membershipOrder project.MembershipOrder
	err := c.ShouldBindQuery(&membershipOrder)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(membershipOrder.GVA_MODEL, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	reMembershipOrder, err := membershipOrderService.GetMembershipOrder(membershipOrder.ID)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}

	response.OkWithData(reMembershipOrder, c)
}

// GetMembershipOrderByOrderNo 根据订单号查询会员订单
// @Tags MembershipOrder
// @Summary 根据订单号查询会员订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param orderNo query string true "订单号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /membershipOrder/findMembershipOrderByOrderNo [get]
func (m *MembershipOrderApi) GetMembershipOrderByOrderNo(c *gin.Context) {
	orderNo := c.Query("orderNo")
	if orderNo == "" {
		response.FailWithMessage("订单号不能为空", c)
		return
	}

	membershipOrder, err := membershipOrderService.GetMembershipOrderByOrderNo(orderNo)
	if err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}

	response.OkWithData(membershipOrder, c)
}

// UpdateMembershipOrderRemark 更新会员订单备注/标记
// @Tags MembershipOrder
// @Summary 更新会员订单备注/标记
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body projectReq.UpdateOrderRemarkReq true "更新会员订单备注"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /membershipOrder/updateMembershipOrderRemark [put]
func (m *MembershipOrderApi) UpdateMembershipOrderRemark(c *gin.Context) {
	var req projectReq.UpdateOrderRemarkReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = membershipOrderService.UpdateMembershipOrderRemark(req)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// CancelMembershipOrder 取消会员订单
// @Tags MembershipOrder
// @Summary 取消会员订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body projectReq.CancelOrderReq true "取消会员订单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"取消成功"}"
// @Router /membershipOrder/cancelMembershipOrder [put]
func (m *MembershipOrderApi) CancelMembershipOrder(c *gin.Context) {
	var req projectReq.CancelOrderReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = membershipOrderService.CancelMembershipOrder(req)
	if err != nil {
		global.GVA_LOG.Error("取消失败!", zap.Error(err))
		response.FailWithMessage("取消失败", c)
		return
	}

	response.OkWithMessage("取消成功", c)
}

// BatchCancelMembershipOrders 批量取消会员订单
// @Tags MembershipOrder
// @Summary 批量取消会员订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量取消会员订单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"取消成功"}"
// @Router /membershipOrder/batchCancelMembershipOrders [put]
func (m *MembershipOrderApi) BatchCancelMembershipOrders(c *gin.Context) {
	var req request.IdsReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = membershipOrderService.BatchCancelMembershipOrders(req.Ids)
	if err != nil {
		global.GVA_LOG.Error("批量取消失败!", zap.Error(err))
		response.FailWithMessage("批量取消失败", c)
		return
	}

	response.OkWithMessage("取消成功", c)
}

// RefundMembershipOrder 申请退款会员订单
// @Tags MembershipOrder
// @Summary 申请退款会员订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body projectReq.RefundOrderReq true "申请退款会员订单 (需要google_auth_code字段)"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"退款申请提交成功"}"
// @Router /membershipOrder/refundMembershipOrder [put]
func (m *MembershipOrderApi) RefundMembershipOrder(c *gin.Context) {
	var req projectReq.RefundOrderReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userID := utils.GetUserID(c)
	userName := utils.GetUserName(c)
	err = membershipOrderRefundService.RefundMembershipOrder(req, &userID, userName)
	if err != nil {
		global.GVA_LOG.Error("退款申请失败!", zap.Error(err))
		response.FailWithMessage("退款申请失败", c)
		return
	}

	response.OkWithMessage("退款申请提交成功", c)
}

// ConfirmPayment 手动确认支付
// @Tags MembershipOrder
// @Summary 手动确认支付
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body projectReq.ConfirmPaymentReq true "手动确认支付 (需要google_auth_code字段)"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"支付确认成功"}"
// @Router /membershipOrder/confirmPayment [put]
func (m *MembershipOrderApi) ConfirmPayment(c *gin.Context) {
	var req projectReq.ConfirmPaymentReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = membershipOrderService.ConfirmPayment(req)
	if err != nil {
		global.GVA_LOG.Error("支付确认失败!", zap.Error(err))
		response.FailWithMessage("支付确认失败", c)
		return
	}

	response.OkWithMessage("支付确认成功", c)
}

// HandlePaymentCallback 处理支付回调
// @Tags MembershipOrder
// @Summary 处理支付回调
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body projectReq.PaymentCallbackReq true "处理支付回调"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"回调处理成功"}"
// @Router /membershipOrder/handlePaymentCallback [post]
func (m *MembershipOrderApi) HandlePaymentCallback(c *gin.Context) {
	var req projectReq.PaymentCallbackReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = membershipOrderService.HandlePaymentCallback(req)
	if err != nil {
		global.GVA_LOG.Error("回调处理失败!", zap.Error(err))
		response.FailWithMessage("回调处理失败", c)
		return
	}

	response.OkWithMessage("回调处理成功", c)
}

// GetOrderStats 获取订单统计信息
// @Tags MembershipOrder
// @Summary 获取订单统计信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query projectReq.OrderStatsReq true "获取订单统计信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /membershipOrder/getOrderStats [get]
func (m *MembershipOrderApi) GetOrderStats(c *gin.Context) {
	var req projectReq.OrderStatsReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	stats, err := membershipOrderService.GetOrderStats(req)
	if err != nil {
		global.GVA_LOG.Error("获取统计信息失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithData(stats, c)
}

// GetUserOrderHistory 获取用户订单历史
// @Tags MembershipOrder
// @Summary 获取用户订单历史
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query projectReq.UserOrderHistoryReq true "获取用户订单历史"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /membershipOrder/getUserOrderHistory [get]
func (m *MembershipOrderApi) GetUserOrderHistory(c *gin.Context) {
	var req projectReq.UserOrderHistoryReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	history, total, err := membershipOrderService.GetUserOrderHistory(req)
	if err != nil {
		global.GVA_LOG.Error("获取用户订单历史失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     history,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// ExportOrders 导出订单数据
// @Tags MembershipOrder
// @Summary 导出订单数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query projectReq.ExportOrderReq true "导出订单数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"导出成功"}"
// @Router /membershipOrder/exportOrders [get]
func (m *MembershipOrderApi) ExportOrders(c *gin.Context) {
	var req projectReq.ExportOrderReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	fileData, fileName, err := membershipOrderService.ExportOrders(req)
	if err != nil {
		global.GVA_LOG.Error("导出失败!", zap.Error(err))
		response.FailWithMessage("导出失败", c)
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(200, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileData)
}

// ValidateOrder 验证订单有效性
// @Tags MembershipOrder
// @Summary 验证订单有效性
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body projectReq.ValidateOrderReq true "验证订单有效性"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"验证成功"}"
// @Router /membershipOrder/validateOrder [post]
func (m *MembershipOrderApi) ValidateOrder(c *gin.Context) {
	var req projectReq.ValidateOrderReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := membershipOrderService.ValidateOrder(req)
	if err != nil {
		global.GVA_LOG.Error("验证失败!", zap.Error(err))
		response.FailWithMessage("验证失败", c)
		return
	}

	response.OkWithData(result, c)
}

// GetPaymentMethods 获取支付方式列表
// @Tags MembershipOrder
// @Summary 获取支付方式列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /membershipOrder/getPaymentMethods [get]
func (m *MembershipOrderApi) GetPaymentMethods(c *gin.Context) {
	methods, err := membershipOrderService.GetPaymentMethods()
	if err != nil {
		global.GVA_LOG.Error("获取支付方式失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithData(methods, c)
}

// GetOrderLogs 获取订单操作日志
// @Tags MembershipOrder
// @Summary 获取订单操作日志
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query projectReq.OrderLogReq true "获取订单操作日志"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /membershipOrder/getOrderLogs [get]
func (m *MembershipOrderApi) GetOrderLogs(c *gin.Context) {
	var req projectReq.OrderLogReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	logs, total, err := membershipOrderService.GetOrderLogs(req)
	if err != nil {
		global.GVA_LOG.Error("获取订单日志失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     logs,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取成功", c)
}

// ManualProcessOrder 手动处理异常订单
// @Tags MembershipOrder
// @Summary 手动处理异常订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body projectReq.ManualProcessOrderReq true "手动处理异常订单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"处理成功"}"
// @Router /membershipOrder/manualProcessOrder [post]
func (m *MembershipOrderApi) ManualProcessOrder(c *gin.Context) {
	var req projectReq.ManualProcessOrderReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = membershipOrderService.ManualProcessOrder(req)
	if err != nil {
		global.GVA_LOG.Error("手动处理失败!", zap.Error(err))
		response.FailWithMessage("处理失败", c)
		return
	}

	response.OkWithMessage("处理成功", c)
}

// GetRefundDetail 获取退款详情
// @Tags MembershipOrder
// @Summary 获取退款详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query projectReq.RefundDetailReq true "获取退款详情"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /membershipOrder/getRefundDetail [get]
func (m *MembershipOrderApi) GetRefundDetail(c *gin.Context) {
	var req projectReq.RefundDetailReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	detail, err := membershipOrderRefundService.GetRefundDetail(req)
	if err != nil {
		global.GVA_LOG.Error("获取退款详情失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithData(detail, c)
}

// SyncPaymentStatus 查询第三方支付状态
// @Tags MembershipOrder
// @Summary 查询第三方支付状态
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body projectReq.QueryPaymentStatusReq true "查询第三方支付状态"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /membershipOrder/syncPaymentStatus [post]
func (m *MembershipOrderApi) SyncPaymentStatus(c *gin.Context) {
	var req projectReq.QueryPaymentStatusReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	result, err := membershipOrderService.SyncPaymentStatus(req)
	if err != nil {
		global.GVA_LOG.Error("同步支付状态失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
		return
	}

	response.OkWithData(result, c)
}

// GetOrderReceipt 获取订单收据
// @Tags MembershipOrder
// @Summary 获取订单收据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query projectReq.OrderReceiptReq true "获取订单收据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /membershipOrder/getOrderReceipt [get]
func (m *MembershipOrderApi) GetOrderReceipt(c *gin.Context) {
	var req projectReq.OrderReceiptReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	receipt, err := membershipOrderService.GetOrderReceipt(req)
	if err != nil {
		global.GVA_LOG.Error("获取订单收据失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}

	response.OkWithData(receipt, c)
}

// SendOrderNotification 发送订单通知
// @Tags MembershipOrder
// @Summary 发送订单通知
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body projectReq.SendOrderNotificationReq true "发送订单通知"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"通知发送成功"}"
// @Router /membershipOrder/sendOrderNotification [post]
func (m *MembershipOrderApi) SendOrderNotification(c *gin.Context) {
	var req projectReq.SendOrderNotificationReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = membershipOrderService.SendOrderNotification(req)
	if err != nil {
		global.GVA_LOG.Error("发送通知失败!", zap.Error(err))
		response.FailWithMessage("通知发送失败", c)
		return
	}

	response.OkWithMessage("通知发送成功", c)
}
