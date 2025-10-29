package project

import (
	"ApkAdmin/constants"
	"ApkAdmin/global"
	projectModel "ApkAdmin/model/project"
	"errors"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	"time"
)

// ==================== 轻量级查询方法 ====================

// GetUserMembershipsForDownload 获取用户会员信息（仅用于下载）
// 只查询必要的数据，不加载用户的其他关联信息
func (u *UserService) GetUserMembershipsForDownload(userID uint) ([]projectModel.UserMembership, error) {
	var memberships []projectModel.UserMembership

	now := time.Now()

	err := global.GVA_DB.
		Where("user_id = ?", userID).
		Where("status = ?", constants.MembershipStatusActive).
		Where("(end_date IS NULL OR end_date > ?)", now).
		Preload("Plan").        // 只加载套餐信息
		Order("end_date DESC"). // 优先返回有效期最长的
		Find(&memberships).Error

	return memberships, err
}

// GetValidMembershipForPlatform 获取指定平台的有效会员
// 返回最佳的可用会员（有效期最长且支持该平台）
func (u *UserService) GetValidMembershipForPlatform(userID uint, platform constants.Platform) (*projectModel.UserMembership, error) {
	var memberships []projectModel.UserMembership

	now := time.Now()

	// 查询所有有效会员
	err := global.GVA_DB.
		Where("user_id = ?", userID).
		Where("status = ?", constants.MembershipStatusActive).
		Where("(end_date IS NULL OR end_date > ?)", now).
		Preload("Plan").
		Order("end_date DESC NULLS FIRST"). // 终身会员优先，然后按有效期排序
		Find(&memberships).Error

	if err != nil {
		return nil, err
	}

	if len(memberships) == 0 {
		return nil, errors.New("用户没有有效会员")
	}

	// 找到支持该平台且可下载的会员
	for i := range memberships {
		if memberships[i].SupportsPlatform(platform.String()) &&
			memberships[i].CanDownload(true, true) {
			return &memberships[i], nil
		}
	}

	// 没有找到合适的会员
	return nil, errors.New("没有支持该平台的可用会员")
}

// CheckUserDownloadPermission 检查用户下载权限（完整检查）
// 返回：会员信息、是否可下载、原因
func (u *UserService) CheckUserDownloadPermission(userID uint, platform constants.Platform) (*projectModel.UserMembership, bool, string) {
	// 1. 获取会员
	membership, err := u.GetValidMembershipForPlatform(userID, platform)
	if err != nil {
		if err.Error() == "用户没有有效会员" {
			return nil, false, "普通用户无法下载，请升级VIP"
		}
		if err.Error() == "没有支持该平台的可用会员" {
			return nil, false, "当前会员不支持该平台"
		}
		return nil, false, "获取会员信息失败"
	}

	// 2. 检查会员状态
	if !membership.IsActive() {
		return nil, false, "会员已失效"
	}

	// 3. 检查下载次数
	if !membership.CanDownload(true, false) {
		return nil, false, "今日下载次数已用完"
	}

	if !membership.CanDownload(false, true) {
		return nil, false, "本月下载次数已用完"
	}

	return membership, true, "success"
}

// ==================== 批量操作方法 ====================

// GetUsersMembershipsMap 批量获取多个用户的会员信息
// 返回 map[userID]memberships
func (u *UserService) GetUsersMembershipsMap(userIDs []uint) (map[uint][]projectModel.UserMembership, error) {
	if len(userIDs) == 0 {
		return make(map[uint][]projectModel.UserMembership), nil
	}

	var memberships []projectModel.UserMembership
	now := time.Now()

	err := global.GVA_DB.
		Where("user_id IN ?", userIDs).
		Where("status = ?", constants.MembershipStatusActive).
		Where("(end_date IS NULL OR end_date > ?)", now).
		Preload("Plan").
		Order("user_id, end_date DESC").
		Find(&memberships).Error

	if err != nil {
		return nil, err
	}

	// 按用户ID分组
	result := make(map[uint][]projectModel.UserMembership)
	for _, membership := range memberships {
		result[membership.UserID] = append(result[membership.UserID], membership)
	}

	return result, nil
}

// ==================== 定时任务 ====================

// UpdateExpiredMemberships 批量更新过期会员状态（定时任务）
// 建议：每小时执行一次
func (u *UserService) UpdateExpiredMemberships() error {
	now := time.Now()

	result := global.GVA_DB.Model(&projectModel.UserMembership{}).
		Where("status = ?", constants.MembershipStatusActive).
		Where("end_date IS NOT NULL").
		Where("end_date < ?", now).
		Update("status", constants.MembershipStatusExpired)

	if result.Error != nil {
		global.GVA_LOG.Error("批量更新过期会员失败", zap.Error(result.Error))
		return result.Error
	}

	if result.RowsAffected > 0 {
		global.GVA_LOG.Info("批量更新过期会员成功",
			zap.Int64("affected", result.RowsAffected))
	}

	return nil
}

// ResetDailyDownloadCount 重置日下载计数（定时任务）
// 建议：每天凌晨执行
func (u *UserService) ResetDailyDownloadCount() error {
	today := time.Now().Truncate(24 * time.Hour)

	result := global.GVA_DB.Model(&projectModel.UserMembership{}).
		Where("last_reset_daily < ?", today).
		Updates(map[string]interface{}{
			"download_used_daily": 0,
			"last_reset_daily":    today,
		})

	if result.Error != nil {
		global.GVA_LOG.Error("重置日下载计数失败", zap.Error(result.Error))
		return result.Error
	}

	global.GVA_LOG.Info("重置日下载计数成功",
		zap.Int64("affected", result.RowsAffected))

	return nil
}

// ResetMonthlyDownloadCount 重置月下载计数（定时任务）
// 建议：每月1号凌晨执行
func (u *UserService) ResetMonthlyDownloadCount() error {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	result := global.GVA_DB.Model(&projectModel.UserMembership{}).
		Where("last_reset_monthly < ?", monthStart).
		Updates(map[string]interface{}{
			"download_used_monthly": 0,
			"last_reset_monthly":    monthStart,
		})

	if result.Error != nil {
		global.GVA_LOG.Error("重置月下载计数失败", zap.Error(result.Error))
		return result.Error
	}

	global.GVA_LOG.Info("重置月下载计数成功",
		zap.Int64("affected", result.RowsAffected))

	return nil
}

// ==================== 统计方法 ====================

// GetMembershipStats 获取会员统计信息
func (u *UserService) GetMembershipStats(userID uint) (map[string]interface{}, error) {
	membership, err := u.GetValidMembershipForPlatform(userID, constants.PlatformAndroid)
	if err != nil {
		return map[string]interface{}{
			"hasMembership":      false,
			"remainingDays":      0,
			"downloadUsedDaily":  0,
			"downloadLimitDaily": 0,
		}, nil
	}

	stats := map[string]interface{}{
		"hasMembership":       true,
		"planName":            membership.PlanName,
		"remainingDays":       membership.GetRemainingDays(),
		"isLifetime":          membership.EndDate == nil,
		"downloadUsedDaily":   membership.DownloadUsedDaily,
		"downloadUsedMonthly": membership.DownloadUsedMonthly,
	}

	// 添加限制信息
	if membership.Plan != nil {
		if membership.Plan.DownloadLimitDaily != nil {
			stats["downloadLimitDaily"] = *membership.Plan.DownloadLimitDaily
			stats["downloadRemainingDaily"] = *membership.Plan.DownloadLimitDaily - int(membership.DownloadUsedDaily)
		}
		if membership.Plan.DownloadLimitMonthly != nil {
			stats["downloadLimitMonthly"] = *membership.Plan.DownloadLimitMonthly
			stats["downloadRemainingMonthly"] = *membership.Plan.DownloadLimitMonthly - int(membership.DownloadUsedMonthly)
		}
	}

	return stats, nil
}

// ==================== 辅助方法 ====================

// IncrementMembershipDownloadCount 增加会员下载计数
func (u *UserService) IncrementMembershipDownloadCount(membershipID uint) error {
	var membership projectModel.UserMembership
	// 使用悲观锁
	err := global.GVA_DB.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", membershipID).
		First(&membership).Error

	if err != nil {
		return err
	}
	// 调用模型方法增加计数
	membership.IncrementDownloadCount()
	// 保存
	return global.GVA_DB.Save(&membership).Error
}

// CheckMembershipDownloadLimit 检查会员下载限制
func (u *UserService) CheckMembershipDownloadLimit(membershipID uint) (canDownload bool, reason string) {
	var membership projectModel.UserMembership

	err := global.GVA_DB.
		Where("id = ?", membershipID).
		Preload("Plan").
		First(&membership).Error

	if err != nil {
		return false, "查询会员信息失败"
	}

	// 检查日限制
	if !membership.CanDownload(true, false) {
		return false, "今日下载次数已用完"
	}

	// 检查月限制
	if !membership.CanDownload(false, true) {
		return false, "本月下载次数已用完"
	}

	return true, "success"
}
