package constants

const (
	SystemName         = "ApkAdmin"            //ApkAdmin
	GoogleVerifyAction = "view-wallet-address" // 谷歌验证器操作
)

// PlanType 套餐类型枚举
type PlanType string

const (
	PlanTypeMonthly  PlanType = "monthly"
	PlanTypeYearly   PlanType = "yearly"
	PlanTypeLifetime PlanType = "lifetime"
)

// OrderStatus 订单状态
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"   // 待支付
	OrderStatusPaid      OrderStatus = "paid"      // 已支付
	OrderStatusFailed    OrderStatus = "failed"    // 支付失败
	OrderStatusRefunded  OrderStatus = "refunded"  // 已退款
	OrderStatusCancelled OrderStatus = "cancelled" // 已取消
)

// ApplicationStatus 应用状态枚举
type ApplicationStatus string

const (
	ApplicationStatusActive    ApplicationStatus = "active"
	ApplicationStatusSuspended ApplicationStatus = "suspended"
	ApplicationStatusDeleted   ApplicationStatus = "deleted"
)

// PackageStatus 包状态枚举
type PackageStatus string

const (
	StatusDraft         PackageStatus = "draft"
	StatusTesting       PackageStatus = "testing"
	StatusReviewPending PackageStatus = "review_pending"
	StatusApproved      PackageStatus = "approved"
	StatusPublished     PackageStatus = "published"
	StatusRejected      PackageStatus = "rejected"
	StatusSuspended     PackageStatus = "suspended"
	StatusArchived      PackageStatus = "archived"
)
