package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	"ApkAdmin/model/project/response"
	"gorm.io/gorm"
	"time"
)

// GetUserDetail  获取用户详情
func (u *UserService) GetUserDetail(conditions ...func(*gorm.DB) *gorm.DB) (*response.UserInfoResponse, error) {

	// 1. 查询用户基本信息
	var user project.User
	query := global.GVA_DB.Model(&project.User{})

	for _, condition := range conditions {
		query = condition(query)
	}
	now := time.Now()
	err := query.Preload("Statistics").
		Preload("CommissionSimple").
		// ✅ 只预加载可用的会员
		Preload("Memberships", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", 1). // 1 = active
								Where("(end_date IS NULL OR end_date > ?)", now).
								Order("created_at DESC")
		}).
		Preload("Memberships.Plan").
		Preload("Referrer", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, username")
		}).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	// 2. 统计直属下级数量
	user.DirectReferralsCount = u.CountDirectReferrals(user.ID)
	// 3. 构建响应
	resp := &response.UserInfoResponse{
		// 基本信息
		ID:       user.ID,
		UUID:     user.UUID.String(),
		Username: user.Username,
		// 联系方式（脱敏）
		Email:     user.MaskEmail(),
		Phone:     user.MaskPhone(),
		EmailFull: user.Email,
		PhoneFull: func() string {
			if user.Phone != nil {
				return *user.Phone
			}
			return ""
		}(),

		// 账户状态
		AccountStatus:     string(user.AccountStatus),
		AccountStatusText: user.GetAccountStatusText(),
		EmailVerified:     user.EmailVerified,
		PhoneVerified:     user.PhoneVerified,
		TwoFactorEnabled:  user.TwoFactorEnabled,

		// 推荐码
		ReferralCode: func() string {
			if user.ReferralCode != nil {
				return *user.ReferralCode
			}
			return ""
		}(),

		// 登录信息
		LoginCount:  user.LoginCount,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,

		// 聚合信息
		IsVerified:    user.IsVerified(),
		HasMembership: user.HasCurrentMembership(),
		IsNewUser:     user.IsNewUser(),
		SecurityLevel: user.GetSecurityLevel(),
	}

	// 4. 佣金信息
	if user.CommissionSimple != nil {
		resp.Commission = &response.CommissionSimple{
			Available: user.CommissionSimple.AvailableAmount,
			Total:     user.CommissionSimple.TotalEarnings,
		}
	}

	// 5. 统计信息
	if user.Statistics != nil {
		resp.Statistics = &response.StatisticsSimple{
			Downloads: user.Statistics.TotalDownloads,
			Orders:    user.Statistics.TotalOrders,
			Spent:     user.Statistics.TotalSpent,
			Referrals: user.Statistics.SuccessfulReferrals,
		}
	}

	// 6. 会员信息
	currentMembership := user.GetCurrentMembership()
	if currentMembership != nil && currentMembership.Plan != nil {
		resp.Membership = &response.MembershipSimple{
			PlanName: currentMembership.Plan.PlanName,
			ExpireAt: currentMembership.EndDate,
		}
	}

	// 7. 推荐信息
	resp.Referral = &response.ReferralInfo{
		DirectCount: user.DirectReferralsCount,
		HasReferrer: user.ReferrerID != nil,
		ReferralCode: func() string {
			if user.ReferralCode != nil {
				return *user.ReferralCode
			}
			return ""
		}(),
	}

	return resp, nil
}

func (u *UserService) CountDirectReferrals(userID uint) int64 {
	var count int64
	global.GVA_DB.Model(&project.User{}).
		Where("referrer_id = ?", userID).
		Where("deleted_at IS NULL"). // 排除软删除
		Count(&count)
	return count
}
