package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	projectReq "ApkAdmin/model/project/request"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type MembershipOrderService struct {
}

// GetMembershipOrderInfoList 分页获取会员订单列表
func (m *MembershipOrderService) GetMembershipOrderInfoList(req projectReq.MembershipOrderSearchRequest) (list []project.MembershipOrder, total int64, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	// 构建查询
	db := global.GVA_DB.Model(&project.MembershipOrder{})

	// 添加搜索条件
	db = m.buildSearchConditions(db, req)

	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	// 获取数据
	err = db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&list).Error
	return list, total, err
}

// buildSearchConditions 构建搜索条件
func (m *MembershipOrderService) buildSearchConditions(db *gorm.DB, req projectReq.MembershipOrderSearchRequest) *gorm.DB {
	// 订单号搜索
	if req.OrderNo != "" {
		db = db.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
	}

	// 用户ID搜索
	if req.UserID != nil && *req.UserID > 0 {
		db = db.Where("user_id = ?", *req.UserID)
	}

	// 套餐类型搜索
	if req.PlanType != "" {
		db = db.Where("plan_type = ?", req.PlanType)
	}

	// 订单类型搜索
	if req.OrderType != "" {
		db = db.Where("order_type = ?", req.OrderType)
	}

	// 订单状态搜索
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	// 购买平台搜索
	if req.Platform != "" {
		db = db.Where("platform = ?", req.Platform)
	}

	// 套餐名称搜索
	if req.PlanName != "" {
		db = db.Where("plan_name LIKE ?", "%"+req.PlanName+"%")
	}

	// 支付方式搜索
	if req.PaymentMethod != "" {
		db = db.Where("payment_method = ?", req.PaymentMethod)
	}

	// 时间范围搜索
	if req.StartTime != "" && req.EndTime != "" {
		db = db.Where("created_at BETWEEN ? AND ?", req.StartTime, req.EndTime)
	} else if len(req.DateRange) == 2 {
		db = db.Where("created_at BETWEEN ? AND ?", req.DateRange[0], req.DateRange[1])
	}

	return db
}

// GetMembershipOrder 根据ID查询会员订单
func (m *MembershipOrderService) GetMembershipOrder(id uint) (membershipOrder project.MembershipOrder, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&membershipOrder).Error
	return
}

// GetMembershipOrderByOrderNo 根据订单号查询会员订单
func (m *MembershipOrderService) GetMembershipOrderByOrderNo(orderNo string) (membershipOrder project.MembershipOrder, err error) {
	err = global.GVA_DB.Where("order_no = ?", orderNo).First(&membershipOrder).Error
	return
}

// UpdateMembershipOrderRemark 更新会员订单备注/标记
func (m *MembershipOrderService) UpdateMembershipOrderRemark(req projectReq.UpdateOrderRemarkReq) error {
	return global.GVA_DB.Model(&project.MembershipOrder{}).
		Where("id = ?", req.ID).
		Update("remark", req.Remark).Error
}

// CancelMembershipOrder 取消会员订单
func (m *MembershipOrderService) CancelMembershipOrder(req projectReq.CancelOrderReq) error {
	// 查询订单
	var order project.MembershipOrder
	err := global.GVA_DB.Where("id = ?", req.ID).First(&order).Error
	if err != nil {
		return err
	}

	// 检查订单状态
	if order.Status != "pending" {
		return errors.New("只能取消待支付订单")
	}

	// 更新订单状态
	return global.GVA_DB.Model(&order).Updates(map[string]interface{}{
		"status":        "cancelled",
		"cancel_reason": req.Reason,
		"updated_at":    time.Now(),
	}).Error
}

// BatchCancelMembershipOrders 批量取消会员订单
func (m *MembershipOrderService) BatchCancelMembershipOrders(ids []int) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 检查所有订单状态
		var count int64
		err := tx.Model(&project.MembershipOrder{}).
			Where("id IN ? AND status = ?", ids, "pending").
			Count(&count).Error
		if err != nil {
			return err
		}

		if int(count) != len(ids) {
			return errors.New("存在非待支付状态的订单，无法批量取消")
		}

		// 批量更新状态
		return tx.Model(&project.MembershipOrder{}).
			Where("id IN ?", ids).
			Updates(map[string]interface{}{
				"status":     "cancelled",
				"updated_at": time.Now(),
			}).Error
	})
}

// ConfirmPayment 手动确认支付
func (m *MembershipOrderService) ConfirmPayment(req projectReq.ConfirmPaymentReq) error {
	// TODO: 验证Google Auth Code
	if !m.validateGoogleAuthCode(req.GoogleAuthCode) {
		return errors.New("Google验证码错误")
	}

	// 查询订单
	var order project.MembershipOrder
	err := global.GVA_DB.Where("id = ?", req.ID).First(&order).Error
	if err != nil {
		return err
	}

	// 检查订单状态
	if order.Status != "pending" {
		return errors.New("订单状态不正确")
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 更新订单状态
		err := tx.Model(&order).Updates(map[string]interface{}{
			"status":       "paid",
			"paid_at":      time.Now(),
			"payment_id":   req.PaymentID,
			"confirm_note": req.Note,
			"updated_at":   time.Now(),
		}).Error
		if err != nil {
			return err
		}

		// TODO: 激活会员服务
		// 这里应该调用会员服务激活相关功能

		return nil
	})
}

// HandlePaymentCallback 处理支付回调
func (m *MembershipOrderService) HandlePaymentCallback(req projectReq.PaymentCallbackReq) error {
	// TODO: 验证回调签名
	if !m.validateCallbackSignature(req) {
		return errors.New("回调签名验证失败")
	}

	// 查询订单
	var order project.MembershipOrder
	err := global.GVA_DB.Where("order_no = ?", req.OrderNo).First(&order).Error
	if err != nil {
		return err
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 根据回调状态更新订单
		updates := map[string]interface{}{
			"payment_id": req.PaymentID,
			"updated_at": time.Now(),
		}

		if req.Status == "success" {
			updates["status"] = "paid"
			updates["paid_at"] = time.Now()
		} else if req.Status == "failed" {
			updates["status"] = "failed"
			updates["fail_reason"] = req.FailReason
		}

		err := tx.Model(&order).Updates(updates).Error
		if err != nil {
			return err
		}

		// 如果支付成功，激活会员服务
		if req.Status == "success" {
			// TODO: 激活会员服务
		}

		return nil
	})
}

// GetOrderStats 获取订单统计信息
func (m *MembershipOrderService) GetOrderStats(req projectReq.OrderStatsReq) (stats projectReq.OrderStatsResp, err error) {
	db := global.GVA_DB.Model(&project.MembershipOrder{})

	// 添加时间范围过滤
	if !req.StartDate.IsZero() && !req.EndDate.IsZero() {
		db = db.Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate)
	}

	// 添加平台过滤
	if req.Platform != "" {
		db = db.Where("platform = ?", req.Platform)
	}

	// 添加套餐类型过滤
	if req.PlanType != "" {
		db = db.Where("plan_type = ?", req.PlanType)
	}

	// 总订单数
	err = db.Count(&stats.TotalOrders).Error
	if err != nil {
		return
	}

	// 各状态订单数
	err = db.Where("status = ?", "paid").Count(&stats.PaidOrders).Error
	if err != nil {
		return
	}

	err = db.Where("status = ?", "pending").Count(&stats.PendingOrders).Error
	if err != nil {
		return
	}

	err = db.Where("status = ?", "cancelled").Count(&stats.CancelledOrders).Error
	if err != nil {
		return
	}

	err = db.Where("status = ?", "refunded").Count(&stats.RefundedOrders).Error
	if err != nil {
		return
	}

	// 总收入
	err = db.Where("status = ?", "paid").Select("COALESCE(SUM(final_amount), 0)").Scan(&stats.TotalRevenue).Error
	if err != nil {
		return
	}

	// 今日统计
	today := time.Now().Format("2006-01-02")
	todayStart := today + " 00:00:00"
	todayEnd := today + " 23:59:59"

	todayDB := global.GVA_DB.Model(&project.MembershipOrder{}).
		Where("created_at BETWEEN ? AND ?", todayStart, todayEnd)

	if req.Platform != "" {
		todayDB = todayDB.Where("platform = ?", req.Platform)
	}

	if req.PlanType != "" {
		todayDB = todayDB.Where("plan_type = ?", req.PlanType)
	}

	err = todayDB.Count(&stats.TodayOrders).Error
	if err != nil {
		return
	}

	err = todayDB.Where("status = ?", "paid").
		Select("COALESCE(SUM(final_amount), 0)").Scan(&stats.TodayRevenue).Error

	return stats, err
}

// GetUserOrderHistory 获取用户订单历史
func (m *MembershipOrderService) GetUserOrderHistory(req projectReq.UserOrderHistoryReq) (history []project.MembershipOrder, total int64, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	db := global.GVA_DB.Model(&project.MembershipOrder{}).Where("user_id = ?", req.UserID)

	// 添加状态过滤
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	// 获取数据
	err = db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&history).Error
	return history, total, err
}

// ExportOrders 导出订单数据
func (m *MembershipOrderService) ExportOrders(req projectReq.ExportOrderReq) (fileData []byte, fileName string, err error) {
	// 查询数据
	var orders []project.MembershipOrder
	db := global.GVA_DB.Model(&project.MembershipOrder{})

	// 添加过滤条件
	if !req.StartDate.IsZero() && !req.EndDate.IsZero() {
		db = db.Where("created_at BETWEEN ? AND ?", req.StartDate, req.EndDate)
	}
	if req.Status != "" {
		db = db.Where("status = ?", req.Status)
	}

	err = db.Find(&orders).Error
	if err != nil {
		return
	}

	// TODO: 生成Excel文件
	// 这里应该使用Excel库生成文件
	fileName = fmt.Sprintf("orders_export_%s.xlsx", time.Now().Format("20060102_150405"))

	return fileData, fileName, nil
}

// ValidateOrder 验证订单有效性
func (m *MembershipOrderService) ValidateOrder(req projectReq.ValidateOrderReq) (result projectReq.ValidateOrderResp, err error) {
	// 查询订单
	var order project.MembershipOrder
	err = global.GVA_DB.Where("order_no = ?", req.OrderNo).First(&order).Error
	if err != nil {
		result.IsValid = false
		result.Message = "订单不存在"
		return result, nil
	}

	// 检查订单状态
	if order.Status == "cancelled" || order.Status == "refunded" {
		result.IsValid = false
		result.Message = "订单已取消或已退款"
		return result, nil
	}

	// 检查订单是否过期
	if time.Now().After(order.ExpiresAt) {
		result.IsValid = false
		result.Message = "订单已过期"
		return result, nil
	}

	result.IsValid = true
	result.Message = "订单有效"
	result.Order = order

	return result, nil
}

// GetPaymentMethods 获取支付方式列表
func (m *MembershipOrderService) GetPaymentMethods() (methods []projectReq.PaymentMethod, err error) {
	// TODO: 从配置或数据库获取支付方式
	methods = []projectReq.PaymentMethod{
		{Code: "alipay", Name: "支付宝", Icon: "alipay.png", Enabled: true},
		{Code: "wechat", Name: "微信支付", Icon: "wechat.png", Enabled: true},
		{Code: "stripe", Name: "信用卡", Icon: "stripe.png", Enabled: true},
		{Code: "paypal", Name: "PayPal", Icon: "paypal.png", Enabled: true},
	}

	return methods, nil
}

// GetOrderLogs 获取订单操作日志
func (m *MembershipOrderService) GetOrderLogs(req projectReq.OrderLogReq) (logs []projectReq.OrderLog, total int64, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	// TODO: 从操作日志表查询
	// 这里应该查询专门的操作日志表
	db := global.GVA_DB.Table("operation_records").
		Where("table_name = ? AND record_id = ?", "membership_orders", req.OrderID)

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&logs).Error
	return logs, total, err
}

// ManualProcessOrder 手动处理异常订单
func (m *MembershipOrderService) ManualProcessOrder(req projectReq.ManualProcessOrderReq) error {
	// 查询订单
	var order project.MembershipOrder
	err := global.GVA_DB.Where("id = ?", req.OrderID).First(&order).Error
	if err != nil {
		return err
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 根据处理类型更新订单
		updates := map[string]interface{}{
			"process_note": req.Note,
			"updated_at":   time.Now(),
		}

		switch req.ProcessType {
		case "confirm_payment":
			updates["status"] = "paid"
			updates["paid_at"] = time.Now()
		case "mark_failed":
			updates["status"] = "failed"
			updates["fail_reason"] = req.Note
		case "force_refund":
			updates["status"] = "refunded"
			updates["refund_reason"] = req.Note
		}

		return tx.Model(&order).Updates(updates).Error
	})
}

// SyncPaymentStatus 查询第三方支付状态
func (m *MembershipOrderService) SyncPaymentStatus(req projectReq.QueryPaymentStatusReq) (result projectReq.PaymentStatusResp, err error) {
	// 查询订单
	var order project.MembershipOrder
	err = global.GVA_DB.Where("order_no = ?", req.OrderNo).First(&order).Error
	if err != nil {
		return
	}

	// TODO: 调用第三方支付接口查询状态
	// 这里应该根据支付方式调用对应的支付接口

	result.OrderNo = order.OrderNo
	result.PaymentStatus = string(order.Status)
	result.PaymentTime = order.PaidAt

	// 如果本地状态与第三方状态不一致，更新本地状态
	if result.PaymentStatus != string(order.Status) {
		global.GVA_DB.Model(&order).Update("status", result.PaymentStatus)
	}

	return result, nil
}

// GetOrderReceipt 获取订单收据
func (m *MembershipOrderService) GetOrderReceipt(req projectReq.OrderReceiptReq) (receipt projectReq.OrderReceiptResp, err error) {
	// 查询订单
	var order project.MembershipOrder
	err = global.GVA_DB.Where("id = ?", req.OrderID).First(&order).Error
	if err != nil {
		return
	}

	if order.Status != "paid" {
		err = errors.New("订单未支付，无法生成收据")
		return
	}

	// 构建收据信息
	receipt.OrderNo = order.OrderNo
	receipt.PlanName = order.PlanName
	receipt.Amount = order.FinalAmount
	receipt.Currency = order.CurrencyCode
	receipt.PaymentMethod = *order.PaymentMethod
	receipt.PaymentTime = *order.PaidAt
	receipt.ReceiptNo = fmt.Sprintf("R%s", order.OrderNo)

	return receipt, nil
}

// SendOrderNotification 发送订单通知
func (m *MembershipOrderService) SendOrderNotification(req projectReq.SendOrderNotificationReq) error {
	// 查询订单
	var order project.MembershipOrder
	err := global.GVA_DB.Where("id = ?", req.OrderID).First(&order).Error
	if err != nil {
		return err
	}

	// TODO: 发送通知
	// 这里应该调用通知服务发送邮件或短信
	switch req.NotificationType {
	case "email":
		// 发送邮件通知
		return m.sendEmailNotification(order, req.Message)
	case "sms":
		// 发送短信通知
		return m.sendSMSNotification(order, req.Message)
	default:
		return errors.New("不支持的通知类型")
	}
}

// validateGoogleAuthCode 验证Google验证码
func (m *MembershipOrderService) validateGoogleAuthCode(code string) bool {
	// TODO: 实现Google Auth验证
	return true // 临时返回true
}

// validateCallbackSignature 验证回调签名
func (m *MembershipOrderService) validateCallbackSignature(req projectReq.PaymentCallbackReq) bool {
	// TODO: 实现签名验证
	return true // 临时返回true
}

// sendEmailNotification 发送邮件通知
func (m *MembershipOrderService) sendEmailNotification(order project.MembershipOrder, message string) error {
	// TODO: 实现邮件发送
	return nil
}

// sendSMSNotification 发送短信通知
func (m *MembershipOrderService) sendSMSNotification(order project.MembershipOrder, message string) error {
	// TODO: 实现短信发送
	return nil
}
