package request

import (
	"ApkAdmin/model/common/request"
	"ApkAdmin/model/project"
	"time"
)

// MembershipOrderSearchRequest 会员订单搜索请求
type MembershipOrderSearchRequest struct {
	request.PageInfo
	OrderNo       string   `json:"order_no" form:"order_no"`             // 订单号
	UserID        *uint    `json:"user_id" form:"user_id"`               // 用户ID
	PlanType      string   `json:"plan_type" form:"plan_type"`           // 套餐类型
	OrderType     string   `json:"order_type" form:"order_type"`         // 订单类型
	Status        string   `json:"status" form:"status"`                 // 订单状态
	Platform      string   `json:"platform" form:"platform"`             // 购买平台
	PlanName      string   `json:"plan_name" form:"plan_name"`           // 套餐名称
	PaymentMethod string   `json:"payment_method" form:"payment_method"` // 支付方式
	StartTime     string   `json:"start_time" form:"start_time"`         // 开始时间
	EndTime       string   `json:"end_time" form:"end_time"`             // 结束时间
	DateRange     []string `json:"date_range" form:"date_range"`         // 时间范围
}

// UpdateOrderRemarkReq 更新订单备注请求
type UpdateOrderRemarkReq struct {
	ID     uint   `json:"id" binding:"required"`    // 订单ID
	Remark string `json:"remark" binding:"max=500"` // 备注内容
}

// CancelOrderReq 取消订单请求
type CancelOrderReq struct {
	ID     uint   `json:"id" binding:"required"`    // 订单ID
	Reason string `json:"reason" binding:"max=500"` // 取消原因
}

// RefundOrderReq 退款订单请求
type RefundOrderReq struct {
	ID             uint     `json:"id" binding:"required"`                     // 订单ID
	RefundReason   string   `json:"refund_reason" binding:"required,min=5"`    // 退款原因
	RefundAmount   *float64 `json:"refund_amount"`                             // 退款金额（可选，为空时退全款）
	RefundType     string   `json:"refund_type"`                               // 退款类型：full-全额退款，partial-部分退款
	GoogleAuthCode string   `json:"google_auth_code" binding:"required,len=6"` // Google验证码
}

// ConfirmPaymentReq 确认支付请求
type ConfirmPaymentReq struct {
	ID             uint   `json:"id" binding:"required"`                     // 订单ID
	PaymentID      string `json:"payment_id"`                                // 支付ID
	Note           string `json:"note"`                                      // 确认备注
	GoogleAuthCode string `json:"google_auth_code" binding:"required,len=6"` // Google验证码
}

// PaymentCallbackReq 支付回调请求
type PaymentCallbackReq struct {
	OrderNo    string `json:"order_no" binding:"required"`   // 订单号
	PaymentID  string `json:"payment_id" binding:"required"` // 支付ID
	Status     string `json:"status" binding:"required"`     // 支付状态
	FailReason string `json:"fail_reason"`                   // 失败原因
	Signature  string `json:"signature" binding:"required"`  // 签名
	Timestamp  int64  `json:"timestamp" binding:"required"`  // 时间戳
	Amount     string `json:"amount"`                        // 金额
}

// OrderStatsReq 订单统计请求
type OrderStatsReq struct {
	StartDate time.Time `json:"start_date" form:"start_date"` // 开始日期
	EndDate   time.Time `json:"end_date" form:"end_date"`     // 结束日期
	Platform  string    `json:"platform" form:"platform"`     // 平台过滤
	PlanType  string    `json:"plan_type" form:"plan_type"`   // 套餐类型过滤
}

// OrderStatsResp 订单统计响应
type OrderStatsResp struct {
	TotalOrders     int64   `json:"total_orders"`     // 总订单数
	PaidOrders      int64   `json:"paid_orders"`      // 已支付订单数
	PendingOrders   int64   `json:"pending_orders"`   // 待支付订单数
	CancelledOrders int64   `json:"cancelled_orders"` // 已取消订单数
	RefundedOrders  int64   `json:"refunded_orders"`  // 已退款订单数
	TotalRevenue    float64 `json:"total_revenue"`    // 总收入
	TodayOrders     int64   `json:"today_orders"`     // 今日订单数
	TodayRevenue    float64 `json:"today_revenue"`    // 今日收入
}

// UserOrderHistoryReq 用户订单历史请求
type UserOrderHistoryReq struct {
	request.PageInfo
	UserID uint   `json:"user_id" form:"user_id" binding:"required"` // 用户ID
	Status string `json:"status" form:"status"`                      // 状态过滤
}

// ExportOrderReq 导出订单请求
type ExportOrderReq struct {
	StartDate     time.Time `json:"start_date" form:"start_date"`         // 开始日期
	EndDate       time.Time `json:"end_date" form:"end_date"`             // 结束日期
	Status        string    `json:"status" form:"status"`                 // 状态过滤
	Platform      string    `json:"platform" form:"platform"`             // 平台过滤
	PlanType      string    `json:"plan_type" form:"plan_type"`           // 套餐类型过滤
	PaymentMethod string    `json:"payment_method" form:"payment_method"` // 支付方式过滤
}

// ValidateOrderReq 验证订单请求
type ValidateOrderReq struct {
	OrderNo string `json:"order_no" binding:"required"` // 订单号
}

// ValidateOrderResp 验证订单响应
type ValidateOrderResp struct {
	IsValid bool          `json:"is_valid"` // 是否有效
	Message string        `json:"message"`  // 验证信息
	Order   project.Order `json:"order"`    // 订单信息
}

// PaymentMethod 支付方式
type PaymentMethod struct {
	Code    string `json:"code"`    // 支付方式代码
	Name    string `json:"name"`    // 支付方式名称
	Icon    string `json:"icon"`    // 图标
	Enabled bool   `json:"enabled"` // 是否启用
}

// OrderLogReq 订单日志请求
type OrderLogReq struct {
	request.PageInfo
	OrderID uint `json:"order_id" form:"order_id" binding:"required"` // 订单ID
}

// OrderLog 订单日志
type OrderLog struct {
	ID          uint      `json:"id"`
	OrderID     uint      `json:"order_id"`
	Action      string    `json:"action"`      // 操作类型
	Description string    `json:"description"` // 操作描述
	OperatorID  uint      `json:"operator_id"` // 操作员ID
	Operator    string    `json:"operator"`    // 操作员名称
	CreatedAt   time.Time `json:"created_at"`  // 创建时间
	IP          string    `json:"ip"`          // 操作IP
	UserAgent   string    `json:"user_agent"`  // 用户代理
}

// ManualProcessOrderReq 手动处理订单请求
type ManualProcessOrderReq struct {
	OrderID     uint   `json:"order_id" binding:"required"`     // 订单ID
	ProcessType string `json:"process_type" binding:"required"` // 处理类型
	Note        string `json:"note" binding:"required"`         // 处理备注
}

// RefundDetailReq 退款详情请求
type RefundDetailReq struct {
	OrderID uint `json:"order_id" form:"order_id" binding:"required"` // 订单ID
}

// RefundDetailResp 退款详情响应
type RefundDetailResp struct {
	ID                 uint       `json:"id"`                    // 退款记录ID
	OrderNo            string     `json:"order_no"`              // 订单号
	RefundAmount       float64    `json:"refund_amount"`         // 退款金额
	RefundStatus       string     `json:"refund_status"`         // 退款状态
	RefundStatusLabel  string     `json:"refund_status_label"`   // 退款状态标签
	RefundType         string     `json:"refund_type"`           // 退款类型
	RefundTypeLabel    string     `json:"refund_type_label"`     // 退款类型标签
	RefundReason       string     `json:"refund_reason"`         // 退款原因
	RefundTime         time.Time  `json:"refund_time"`           // 申请退款时间
	ProcessedAt        *time.Time `json:"processed_at"`          // 处理时间
	CompletedAt        *time.Time `json:"completed_at"`          // 完成时间
	ThirdPartyRefundID string     `json:"third_party_refund_id"` // 第三方退款ID
	OperatorName       string     `json:"operator_name"`         // 操作员名称
	FailureReason      string     `json:"failure_reason"`        // 失败原因
}

// RefundListReq 退款记录列表请求
type RefundListReq struct {
	request.PageInfo
	OrderNo      string `json:"order_no" form:"order_no"`           // 订单号
	RefundStatus string `json:"refund_status" form:"refund_status"` // 退款状态
	RefundType   string `json:"refund_type" form:"refund_type"`     // 退款类型
	StartTime    string `json:"start_time" form:"start_time"`       // 开始时间
	EndTime      string `json:"end_time" form:"end_time"`           // 结束时间
}

// QueryPaymentStatusReq 查询支付状态请求
type QueryPaymentStatusReq struct {
	OrderNo string `json:"order_no" binding:"required"` // 订单号
}

// PaymentStatusResp 支付状态响应
type PaymentStatusResp struct {
	OrderNo       string     `json:"order_no"`       // 订单号
	PaymentStatus string     `json:"payment_status"` // 支付状态
	PaymentTime   *time.Time `json:"payment_time"`   // 支付时间
	PaymentID     string     `json:"payment_id"`     // 支付ID
	Amount        string     `json:"amount"`         // 支付金额
	ThirdStatus   string     `json:"third_status"`   // 第三方状态
}

// OrderReceiptReq 订单收据请求
type OrderReceiptReq struct {
	OrderID uint `json:"order_id" form:"order_id" binding:"required"` // 订单ID
}

// OrderReceiptResp 订单收据响应
type OrderReceiptResp struct {
	OrderNo       string    `json:"order_no"`       // 订单号
	ReceiptNo     string    `json:"receipt_no"`     // 收据号
	PlanName      string    `json:"plan_name"`      // 套餐名称
	Amount        float64   `json:"amount"`         // 金额
	Currency      string    `json:"currency"`       // 货币
	PaymentMethod string    `json:"payment_method"` // 支付方式
	PaymentTime   time.Time `json:"payment_time"`   // 支付时间
	CompanyName   string    `json:"company_name"`   // 公司名称
	CompanyAddr   string    `json:"company_addr"`   // 公司地址
}

// SendOrderNotificationReq 发送订单通知请求
type SendOrderNotificationReq struct {
	OrderID          uint     `json:"order_id" binding:"required"`          // 订单ID
	NotificationType string   `json:"notification_type" binding:"required"` // 通知类型
	Message          string   `json:"message"`                              // 自定义消息
	Recipients       []string `json:"recipients"`                           // 接收人列表
}
