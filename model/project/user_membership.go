package project

import (
	"ApkAdmin/constants"
	"encoding/json"
	"time"
)

// UserMembership 用户会员记录表
type UserMembership struct {
	ID                  uint                       `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	UserID              uint                       `json:"user_id" gorm:"not null;index:idx_user_id;comment:用户ID"`
	OrderID             *uint                      `json:"order_id" gorm:"index:idx_order_id;comment:关联订单ID"`
	PlanID              uint                       `json:"plan_id" gorm:"not null;index:idx_plan_id;comment:套餐ID"`
	PlanCode            string                     `json:"plan_code" gorm:"type:varchar(50);not null;comment:套餐代码快照"`
	PlanName            string                     `json:"plan_name" gorm:"type:varchar(100);not null;comment:套餐名称快照"`
	Detail              string                     `json:"detail" gorm:"type:text;comment:权益详情"`
	Status              constants.MembershipStatus `json:"status" gorm:"not null;default:1;index:idx_status;comment:1生效中,2已过期,3已取消,4已暂停,5已被替代（升级/降级)"`
	StartDate           time.Time                  `json:"start_date" gorm:"not null;comment:开始时间"`
	EndDate             *time.Time                 `json:"end_date" gorm:"index:idx_end_date;comment:结束时间"`
	AutoRenew           bool                       `json:"auto_renew" gorm:"default:0;comment:是否自动续费"`
	DownloadUsedDaily   uint                       `json:"download_used_daily" gorm:"default:0;comment:今日已用下载次数"`
	DownloadUsedMonthly uint                       `json:"download_used_monthly" gorm:"default:0;comment:本月已用下载次数"`
	LastResetDaily      *time.Time                 `json:"last_reset_daily" gorm:"type:date;comment:上次重置日下载计数的日期"`
	LastResetMonthly    *time.Time                 `json:"last_reset_monthly" gorm:"type:date;comment:上次重置月下载计数的日期"`
	ReplacedBy          *uint                      `json:"replaced_by" gorm:"index:idx_replaced_by;comment:被哪个会员记录替代"`
	CreatedAt           time.Time                  `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt           time.Time                  `json:"updated_at" gorm:"comment:更新时间"`

	// 关联关系
	Plan             *MembershipPlan `json:"plan,omitempty" gorm:"foreignKey:PlanID"`
	Order            *Order          `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	ReplacedByRecord *UserMembership `json:"replaced_by_record,omitempty" gorm:"foreignKey:ReplacedBy"`
}

// TableName 指定表名
func (UserMembership) TableName() string {
	return "user_memberships"
}

// IsActive 检查会员是否有效
func (u *UserMembership) IsActive() bool {
	// 首先检查状态是否为 active
	if u.Status.IsActive() {
		return false
	}
	// 如果没有结束时间，说明是终身会员，直接返回 true
	if u.EndDate == nil {
		return true
	}
	// 检查是否已过期
	return time.Now().Before(*u.EndDate)
}

// IsExpired 检查是否已过期
func (u *UserMembership) IsExpired() bool {
	if u.EndDate == nil {
		return false // 终身会员不过期
	}
	return time.Now().After(*u.EndDate)
}

// CanDownload 检查是否可以下载（考虑下载限制）
func (u *UserMembership) CanDownload(checkDaily, checkMonthly bool) bool {
	if !u.IsActive() {
		return false
	}
	// 获取套餐限制（需要加载Plan关联）
	if u.Plan == nil {
		return false
	}
	// 检查日下载限制
	if checkDaily && u.Plan.DownloadLimitDaily != nil {
		if u.DownloadUsedDaily >= uint(*u.Plan.DownloadLimitDaily) {
			return false
		}
	}
	// 检查月下载限制
	if checkMonthly && u.Plan.DownloadLimitMonthly != nil {
		if u.DownloadUsedMonthly >= uint(*u.Plan.DownloadLimitMonthly) {
			return false
		}
		return true
	}
	return false
}

// IncrementDownloadCount 增加下载计数
func (u *UserMembership) IncrementDownloadCount() {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 检查是否需要重置日计数
	if u.LastResetDaily == nil || u.LastResetDaily.Before(today) {
		u.DownloadUsedDaily = 0
		u.LastResetDaily = &today
	}

	// 检查是否需要重置月计数
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	if u.LastResetMonthly == nil || u.LastResetMonthly.Before(monthStart) {
		u.DownloadUsedMonthly = 0
		u.LastResetMonthly = &monthStart
	}

	// 增加计数
	u.DownloadUsedDaily++
	u.DownloadUsedMonthly++
}

// GetRemainingDays 获取剩余天数
func (u *UserMembership) GetRemainingDays() int {
	if u.EndDate == nil {
		return -1 // 终身会员返回 -1
	}

	if !u.IsActive() {
		return 0
	}

	remaining := time.Until(*u.EndDate)
	return int(remaining.Hours() / 24)
}

// SupportsPlatform 检查当前会员套餐是否支持指定平台
func (u *UserMembership) SupportsPlatform(platform string) bool {
	if u.Plan == nil {
		return false
	}

	// 解析套餐支持的平台列表
	var platforms []string
	if err := json.Unmarshal(u.Plan.Platform, &platforms); err != nil {
		return false
	}

	// 检查目标平台是否在支持列表中
	for _, p := range platforms {
		if p == platform {
			return true
		}
	}
	return false
}
