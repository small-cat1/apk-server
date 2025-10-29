package response

import "time"

// UserInfoResponse 用户信息响应（方案1：最小改动）
type UserInfoResponse struct {
	// 基本信息
	ID       uint   `json:"id"`
	UUID     string `json:"uuid"`
	Username string `json:"username"`

	// 联系方式（脱敏）
	Email     string `json:"email"`     // 脱敏邮箱
	Phone     string `json:"phone"`     // 脱敏手机
	EmailFull string `json:"emailFull"` // 完整邮箱
	PhoneFull string `json:"phoneFull"` // 完整手机

	// 账户状态
	AccountStatus     string `json:"accountStatus"`
	AccountStatusText string `json:"accountStatusText"`
	EmailVerified     bool   `json:"emailVerified"`
	PhoneVerified     bool   `json:"phoneVerified"`
	TwoFactorEnabled  bool   `json:"twoFactorEnabled"`

	// 推荐码
	ReferralCode string `json:"referralCode,omitempty"`

	// 登录信息
	LoginCount  uint       `json:"loginCount"`
	LastLoginAt *time.Time `json:"lastLoginAt,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`

	// 聚合信息
	IsVerified    bool   `json:"isVerified"`
	HasMembership bool   `json:"hasMembership"`
	IsNewUser     bool   `json:"isNewUser"`
	SecurityLevel string `json:"securityLevel"`

	// 佣金信息（简化）
	Commission *CommissionSimple `json:"commission,omitempty"`

	// 统计信息（简化）
	Statistics *StatisticsSimple `json:"statistics,omitempty"`

	// 会员信息
	Membership *MembershipSimple `json:"membership,omitempty"`

	// 推荐信息
	Referral *ReferralInfo `json:"referral,omitempty"`
}

// CommissionSimple 佣金信息（简化）
type CommissionSimple struct {
	Available float64 `json:"available"` // 可用金额
	Total     float64 `json:"total"`     // 累计收益
}

// StatisticsSimple 统计信息（简化）
type StatisticsSimple struct {
	Downloads uint    `json:"downloads"`
	Orders    uint    `json:"orders"`
	Spent     float64 `json:"spent"`
	Referrals uint    `json:"referrals"`
}

// MembershipSimple 会员信息（简化）
type MembershipSimple struct {
	PlanName string     `json:"planName"`
	ExpireAt *time.Time `json:"expireAt,omitempty"`
}

// ReferralInfo 推荐信息
type ReferralInfo struct {
	DirectCount  int64  `json:"directCount"`
	HasReferrer  bool   `json:"hasReferrer"`
	ReferralCode string `json:"referralCode,omitempty"`
}
