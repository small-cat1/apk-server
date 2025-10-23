package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	projectReq "ApkAdmin/model/project/request"
	"errors"
	"time"

	"gorm.io/gorm"
)

type MembershipOrderRefundService struct {
}

// RefundMembershipOrder 申请退款会员订单
func (r *MembershipOrderRefundService) RefundMembershipOrder(req projectReq.RefundOrderReq, OperatorID *uint, OperatorName string) error {
	// TODO: 验证Google Auth Code
	if !r.validateGoogleAuthCode(req.GoogleAuthCode) {
		return errors.New("Google验证码错误")
	}

	// 查询订单
	var order project.Order
	err := global.GVA_DB.Where("id = ?", req.ID).First(&order).Error
	if err != nil {
		return err
	}

	// 检查订单状态
	if order.Status != "paid" {
		return errors.New("只能对已支付订单申请退款")
	}

	// 检查是否已有待处理的退款申请
	var existingRefund project.MembershipOrderRefund
	err = global.GVA_DB.Where("order_id = ? AND refund_status IN ?", order.ID, []string{"pending", "processing"}).First(&existingRefund).Error
	if err == nil {
		return errors.New("该订单已有待处理的退款申请")
	}

	// 确定退款金额
	refundAmount := order.FinalAmount // 默认全额退款
	refundType := "full"

	if req.RefundAmount != nil {
		// 验证部分退款金额
		if *req.RefundAmount <= 0 || *req.RefundAmount > order.FinalAmount {
			return errors.New("退款金额不能超过订单实际支付金额")
		}
		refundAmount = *req.RefundAmount
		refundType = "partial"
	}

	if req.RefundType != "" {
		refundType = req.RefundType
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 创建退款记录
		refund := project.MembershipOrderRefund{
			OrderID:      uint(order.ID),
			OrderNo:      order.OrderNo,
			RefundAmount: refundAmount,
			RefundReason: req.RefundReason,
			RefundType:   refundType,
			RefundStatus: "pending",
			// TODO: 设置操作员信息
			OperatorID:   OperatorID,
			OperatorName: OperatorName,
		}

		err := tx.Create(&refund).Error
		if err != nil {
			return err
		}

		// 更新订单状态为已退款
		err = tx.Model(&order).Updates(map[string]interface{}{
			"status":     "refunded",
			"updated_at": time.Now(),
		}).Error
		if err != nil {
			return err
		}

		// TODO: 调用第三方支付接口进行退款
		// 这里应该调用具体的支付服务进行退款处理
		// 成功后调用 ProcessRefund 方法更新退款记录状态和第三方退款ID

		return nil
	})
}

// GetRefundDetail 获取退款详情
func (r *MembershipOrderRefundService) GetRefundDetail(req projectReq.RefundDetailReq) (detail projectReq.RefundDetailResp, err error) {
	// 查询订单
	var order project.Order
	err = global.GVA_DB.Where("id = ?", req.OrderID).First(&order).Error
	if err != nil {
		return
	}

	if order.Status != "refunded" {
		err = errors.New("订单未退款")
		return
	}

	// 查询退款记录
	var refund project.MembershipOrderRefund
	err = global.GVA_DB.Where("order_id = ?", req.OrderID).Order("created_at DESC").First(&refund).Error
	if err != nil {
		return
	}

	// 构建响应
	detail = projectReq.RefundDetailResp{
		ID:                 refund.ID,
		OrderNo:            refund.OrderNo,
		RefundAmount:       refund.RefundAmount,
		RefundStatus:       refund.RefundStatus,
		RefundStatusLabel:  refund.GetRefundStatusLabel(),
		RefundType:         refund.RefundType,
		RefundTypeLabel:    refund.GetRefundTypeLabel(),
		RefundReason:       refund.RefundReason,
		RefundTime:         refund.CreatedAt,
		ProcessedAt:        refund.ProcessedAt,
		CompletedAt:        refund.CompletedAt,
		ThirdPartyRefundID: refund.ThirdPartyRefundID,
		OperatorName:       refund.OperatorName,
		FailureReason:      refund.FailureReason,
	}

	return detail, nil
}

// GetRefundList 获取退款记录列表
func (r *MembershipOrderRefundService) GetRefundList(req projectReq.RefundListReq) (list []project.MembershipOrderRefund, total int64, err error) {
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	// 构建查询
	db := global.GVA_DB.Model(&project.MembershipOrderRefund{})

	// 添加搜索条件
	if req.OrderNo != "" {
		db = db.Where("order_no LIKE ?", "%"+req.OrderNo+"%")
	}

	if req.RefundStatus != "" {
		db = db.Where("refund_status = ?", req.RefundStatus)
	}

	if req.RefundType != "" {
		db = db.Where("refund_type = ?", req.RefundType)
	}

	// 时间范围搜索
	if req.StartTime != "" && req.EndTime != "" {
		db = db.Where("created_at BETWEEN ? AND ?", req.StartTime, req.EndTime)
	}

	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	// 获取数据，预加载订单信息
	err = db.Preload("Order").Limit(limit).Offset(offset).Order("created_at DESC").Find(&list).Error
	return list, total, err
}

// ProcessRefund 处理退款（更新退款状态）
func (r *MembershipOrderRefundService) ProcessRefund(refundID uint, status string, thirdPartyRefundID string, failureReason string) error {
	var refund project.MembershipOrderRefund
	err := global.GVA_DB.Where("id = ?", refundID).First(&refund).Error
	if err != nil {
		return err
	}

	// 检查状态流转是否合法
	if !r.isValidStatusTransition(refund.RefundStatus, status) {
		return errors.New("无效的状态转换")
	}

	updates := map[string]interface{}{
		"refund_status": status,
		"updated_at":    time.Now(),
	}

	// 设置处理时间
	if status == "processing" && refund.ProcessedAt == nil {
		now := time.Now()
		updates["processed_at"] = &now
	}

	// 设置完成时间和第三方退款ID
	if status == "success" {
		now := time.Now()
		updates["completed_at"] = &now
		if thirdPartyRefundID != "" {
			updates["third_party_refund_id"] = thirdPartyRefundID
		}
	}

	// 设置失败原因
	if status == "failed" && failureReason != "" {
		updates["failure_reason"] = failureReason
	}

	return global.GVA_DB.Model(&refund).Updates(updates).Error
}

// CancelRefund 取消退款申请
func (r *MembershipOrderRefundService) CancelRefund(refundID uint, reason string) error {
	var refund project.MembershipOrderRefund
	err := global.GVA_DB.Where("id = ?", refundID).First(&refund).Error
	if err != nil {
		return err
	}

	if !refund.CanCancel() {
		return errors.New("当前状态不允许取消")
	}

	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 更新退款记录状态
		err := tx.Model(&refund).Updates(map[string]interface{}{
			"refund_status":  "cancelled",
			"failure_reason": reason,
			"updated_at":     time.Now(),
		}).Error
		if err != nil {
			return err
		}

		// 恢复订单状态为已支付
		return tx.Model(&project.Order{}).
			Where("id = ?", refund.OrderID).
			Update("status", "paid").Error
	})
}

// RetryRefund 重试退款
func (r *MembershipOrderRefundService) RetryRefund(refundID uint) error {
	var refund project.MembershipOrderRefund
	err := global.GVA_DB.Where("id = ?", refundID).First(&refund).Error
	if err != nil {
		return err
	}

	if !refund.CanRetry() {
		return errors.New("当前状态不允许重试")
	}

	// 重置状态为待处理
	updates := map[string]interface{}{
		"refund_status":  "pending",
		"failure_reason": "",
		"updated_at":     time.Now(),
	}

	err = global.GVA_DB.Model(&refund).Updates(updates).Error
	if err != nil {
		return err
	}

	// TODO: 重新调用第三方支付接口进行退款

	return nil
}

// GetRefundByOrderID 根据订单ID获取退款记录
func (r *MembershipOrderRefundService) GetRefundByOrderID(orderID uint) (refunds []project.MembershipOrderRefund, err error) {
	err = global.GVA_DB.Where("order_id = ?", orderID).Order("created_at DESC").Find(&refunds).Error
	return
}

// GetRefundByID 根据退款ID获取退款记录
func (r *MembershipOrderRefundService) GetRefundByID(refundID uint) (refund project.MembershipOrderRefund, err error) {
	err = global.GVA_DB.Preload("Order").Where("id = ?", refundID).First(&refund).Error
	return
}

// GetRefundStats 获取退款统计信息
func (r *MembershipOrderRefundService) GetRefundStats(startTime, endTime time.Time) (stats map[string]interface{}, err error) {
	db := global.GVA_DB.Model(&project.MembershipOrderRefund{})

	if !startTime.IsZero() && !endTime.IsZero() {
		db = db.Where("created_at BETWEEN ? AND ?", startTime, endTime)
	}

	stats = make(map[string]interface{})

	// 总退款申请数
	var totalCount int64
	db.Count(&totalCount)
	stats["total_refunds"] = totalCount

	// 各状态统计
	statusStats := make(map[string]int64)
	for _, status := range project.RefundStatusOptions {
		var count int64
		db.Where("refund_status = ?", status).Count(&count)
		statusStats[status] = count
	}
	stats["status_stats"] = statusStats

	// 退款金额统计
	var totalAmount float64
	db.Where("refund_status = ?", "success").Select("COALESCE(SUM(refund_amount), 0)").Scan(&totalAmount)
	stats["total_refund_amount"] = totalAmount

	return stats, nil
}

// validateGoogleAuthCode 验证Google验证码
func (r *MembershipOrderRefundService) validateGoogleAuthCode(code string) bool {
	// TODO: 实现Google Auth验证
	return true // 临时返回true
}

// isValidStatusTransition 检查状态转换是否有效
func (r *MembershipOrderRefundService) isValidStatusTransition(fromStatus, toStatus string) bool {
	// 定义有效的状态转换规则
	validTransitions := map[string][]string{
		"pending":    {"processing", "cancelled"},
		"processing": {"success", "failed"},
		"failed":     {"pending"}, // 允许重试
		"success":    {},          // 成功状态不能转换
		"cancelled":  {"pending"}, // 取消后可以重新申请
	}

	allowedNext, exists := validTransitions[fromStatus]
	if !exists {
		return false
	}

	for _, allowed := range allowedNext {
		if allowed == toStatus {
			return true
		}
	}

	return false
}
