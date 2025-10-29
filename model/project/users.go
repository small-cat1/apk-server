package project

import (
	"ApkAdmin/constants"
	"encoding/json"
	"github.com/google/uuid"
	"strings"
	"time"
)

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

	// ✅ 新增这一行
	DirectReferralsCount int64 `json:"directReferralsCount" gorm:"-"`
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

// MaskPhone 脱敏手机号
func (u *User) MaskPhone() string {
	if u.Phone == nil || *u.Phone == "" {
		return ""
	}
	phone := *u.Phone
	if len(phone) == 11 {
		return phone[:3] + "****" + phone[7:]
	}
	return phone
}

// MaskEmail 脱敏邮箱
func (u *User) MaskEmail() string {
	if u.Email == "" {
		return ""
	}

	parts := strings.Split(u.Email, "@")
	if len(parts) != 2 {
		return u.Email
	}

	username := parts[0]
	domain := parts[1]

	if len(username) <= 1 {
		return "*@" + domain
	} else if len(username) <= 3 {
		return username[:1] + "**@" + domain
	} else {
		return username[:1] + "***@" + domain
	}
}

// GetAccountStatusText 获取账户状态文本
func (u *User) GetAccountStatusText() string {
	switch u.AccountStatus {
	case constants.AccountStatusNormal:
		return "正常"
	case constants.AccountStatusSuspended:
		return "已暂停"
	case constants.AccountStatusDeleted:
		return "已删除"
	case constants.AccountStatusPending:
		return "待验证"
	case constants.AccountStatusDisabled:
		return "已封禁"
	default:
		return "未知"
	}
}

// IsNewUser 是否新用户（注册7天内）
func (u *User) IsNewUser() bool {
	return time.Since(u.CreatedAt) <= 7*24*time.Hour
}

// GetSecurityLevel 获取安全等级
func (u *User) GetSecurityLevel() string {
	score := u.GetSecurityScore()
	if score >= 80 {
		return "high"
	} else if score >= 50 {
		return "medium"
	}
	return "low"
}

// GetSecurityScore 计算安全分数
func (u *User) GetSecurityScore() int {
	score := 0

	if u.EmailVerified {
		score += 30
	}
	if u.PhoneVerified {
		score += 30
	}
	if u.TwoFactorEnabled {
		score += 40
	}

	return score
}

func (User) TableName() string {
	return "users"
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
