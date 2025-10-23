package project

import (
	"ApkAdmin/constants"
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// AccountStatus 账户状态
type AccountStatus string

const (
	AccountStatusActive              AccountStatus = "active"               // 正常
	AccountStatusSuspended           AccountStatus = "suspended"            // 暂停
	AccountStatusDeleted             AccountStatus = "deleted"              // 已删除
	AccountStatusPendingVerification AccountStatus = "pending_verification" // 待验证
	AccountStatusBanned              AccountStatus = "banned"               // 被封禁
)

func (s AccountStatus) IsNormal() bool {
	return s == AccountStatusActive
}

type UserLogin interface {
	GetUsername() string
	GetUUID() uuid.UUID
	GetUserId() uint
	GetUserInfo() any
	GetEmail() string
}

// User 用户表
type User struct {
	ID               uint                    `json:"id" gorm:"primaryKey;autoIncrement;comment:用户ID"`
	UUID             uuid.UUID               `json:"uuid" gorm:"type:varchar(36);not null;uniqueIndex:uk_uuid;comment:用户UUID"`
	Username         string                  `json:"username" gorm:"type:varchar(50);not null;uniqueIndex:uk_username;comment:用户名"`
	Email            string                  `json:"email" gorm:"type:varchar(100);not null;uniqueIndex:uk_email;comment:邮箱"`
	Phone            *string                 `json:"phone" gorm:"type:varchar(20);index:idx_phone;comment:用户手机号"`
	PasswordHash     string                  `json:"-" gorm:"type:varchar(255);not null;comment:密码哈希"` // 不返回给前端
	AccountStatus    constants.AccountStatus `json:"account_status" gorm:"default:1;index:idx_account_status;comment:账户状态"`
	StatusReason     string                  `json:"status_reason" gorm:"size:255;comment:状态原因"`
	StatusExpireAt   *time.Time              `json:"status_expire_at" gorm:"comment:状态过期时间"`
	EmailVerified    bool                    `json:"email_verified" gorm:"default:0;comment:邮箱是否已验证"`
	PhoneVerified    bool                    `json:"phone_verified" gorm:"default:0;comment:手机是否已验证"`
	TwoFactorEnabled bool                    `json:"two_factor_enabled" gorm:"default:0;comment:是否启用双因子认证"`
	RegisterIP       *string                 `json:"register_ip" gorm:"type:varchar(45);comment:注册IP"`
	// 登录成功记录
	LastLoginAt     *time.Time `json:"last_login_at" gorm:"index:idx_last_login;comment:最后登录时间"`
	LastLoginIP     *string    `json:"last_login_ip" gorm:"type:varchar(45);comment:最后登录IP"`
	LastLoginDevice *string    `json:"last_login_device" gorm:"type:varchar(255);comment:最后登录设备"`
	LoginCount      uint       `json:"login_count" gorm:"default:0;comment:登录次数"`
	// ⭐ 登录失败记录（新增）
	FailedLoginAttempts uint       `json:"failed_login_attempts" gorm:"default:0;comment:连续失败登录次数"`
	LastFailedLoginAt   *time.Time `json:"last_failed_login_at" gorm:"comment:最后失败登录时间"`
	LastFailedLoginIP   string     `json:"last_failed_login_ip" gorm:"size:45;comment:最后失败登录IP"`

	// 推荐人
	ReferrerID   *uint   `json:"referrer_id" gorm:"index:idx_referrer_id;comment:推荐人ID"`
	ReferralCode *string `json:"referral_code" gorm:"type:varchar(20);uniqueIndex:uk_referral_code;comment:专属推荐码"`

	DownloadPreferences     json.RawMessage `json:"download_preferences" gorm:"type:json;comment:下载偏好设置"`
	NotificationPreferences json.RawMessage `json:"notification_preferences" gorm:"type:json;comment:通知偏好设置"`
	PrivacySettings         json.RawMessage `json:"privacy_settings" gorm:"type:json;comment:隐私设置"`
	CreatedAt               time.Time       `json:"created_at" gorm:"index:idx_created_at;comment:创建时间"`
	UpdatedAt               time.Time       `json:"updated_at" gorm:"comment:更新时间"`
	DeletedAt               *time.Time      `json:"deleted_at" gorm:"comment:软删除时间"`

	// 关联关系
	Commission       *UserCommissionAccount       `json:"commission,omitempty" gorm:"foreignKey:UserID"`
	CommissionSimple *UserCommissionAccountSimple `json:"commissionSimple,omitempty" gorm:"foreignKey:UserID"` // 新增
	Statistics       *UserStatistics              `json:"statistics,omitempty" gorm:"foreignKey:UserID"`
	Memberships      []UserMembership             `json:"memberships" gorm:"foreignKey:UserID"`
	Referrer         *User                        `json:"referrer,omitempty" gorm:"foreignKey:ReferrerID"`
	Referrals        []User                       `json:"referrals,omitempty" gorm:"foreignKey:ReferrerID"`
}

func (u User) GetUsername() string {
	return u.Username
}

func (u User) GetUUID() uuid.UUID {
	return u.UUID
}

func (u User) GetUserId() uint {
	return u.ID
}

func (u User) GetUserInfo() any {
	return u
}

func (u User) GetEmail() string {
	return u.Email
}

// UserStatistics 用户统计表
type UserStatistics struct {
	UserID              uint       `json:"user_id" gorm:"primaryKey;comment:用户ID"`
	TotalDownloads      uint       `json:"total_downloads" gorm:"default:0;comment:总下载次数"`
	TotalSpent          float64    `json:"total_spent" gorm:"type:decimal(10,2);default:0.00;comment:总消费金额"`
	TotalOrders         uint       `json:"total_orders" gorm:"default:0;comment:总订单数"`
	SuccessfulReferrals uint       `json:"successful_referrals" gorm:"default:0;comment:成功推荐人数"`
	LastDownloadAt      *time.Time `json:"last_download_at" gorm:"comment:最后下载时间"`
	LastOrderAt         *time.Time `json:"last_order_at" gorm:"comment:最后订单时间"`
	CreatedAt           time.Time  `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt           time.Time  `json:"updated_at" gorm:"comment:更新时间"`

	// 关联关系
	User *User `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (User) TableName() string {
	return "users"
}

func (UserStatistics) TableName() string {
	return "user_statistics"
}

// IsVerified 检查用户是否已完成验证
func (u *User) IsVerified() bool {
	return u.EmailVerified && (u.Phone == nil || u.PhoneVerified)
}

// HasCurrentMembership 检查用户是否有有效会员
func (u *User) HasCurrentMembership() bool {
	if u.Memberships == nil {
		return false
	}

	for _, membership := range u.Memberships {
		if membership.IsActive() {
			return true
		}
	}
	return false
}

// GetCurrentMembership 获取用户当前有效会员
func (u *User) GetCurrentMembership() *UserMembership {
	if u.Memberships == nil {
		return nil
	}

	for _, membership := range u.Memberships {
		if membership.IsActive() {
			return &membership
		}
	}
	return nil
}
