package project

import (
	"ApkAdmin/constants"
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	"ApkAdmin/model/project/request"
	"ApkAdmin/model/project/response"
	"ApkAdmin/utils"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type UserService struct{}

// GetUserList 获取用户列表
func (u *UserService) GetUserList(info request.UserListRequest, order string, desc bool) (list interface{}, total int64, err error) {
	var userList []project.User
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	// 构建查询条件
	db := global.GVA_DB.Model(&project.User{})
	db = u.buildSearchConditions(db, info)

	// 预加载关联数据
	db = db.Preload("Statistics").Preload("Memberships", func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", "active").Order("created_at DESC")
	})

	// 获取总数
	err = db.Count(&total).Error
	if err != nil {
		return userList, total, err
	}

	// 分页和排序
	db = db.Limit(limit).Offset(offset)
	orderStr := u.buildOrderConditions(order, desc)

	err = db.Order(orderStr).Find(&userList).Error

	// 为每个用户添加当前会员信息
	for i := range userList {
		if len(userList[i].Memberships) > 0 {
			// 添加一个临时字段用于前端显示
			// 这里可以通过自定义结构体或者map来实现
		}
	}

	return userList, total, err
}

// buildSearchConditions 构建搜索条件
func (u *UserService) buildSearchConditions(db *gorm.DB, info request.UserListRequest) *gorm.DB {
	// 用户名搜索（模糊搜索）
	if info.Username != "" {
		db = db.Where("username LIKE ?", "%"+info.Username+"%")
	}

	// 邮箱搜索（模糊搜索）
	if info.Email != "" {
		db = db.Where("email LIKE ?", "%"+info.Email+"%")
	}

	// 手机号搜索（模糊搜索）
	if info.Phone != "" {
		db = db.Where("phone LIKE ?", "%"+info.Phone+"%")
	}

	// 账户状态搜索
	if info.AccountStatus.IsValid() {
		db = db.Where("account_status = ?", info.AccountStatus)
	}

	// 性别搜索
	if info.Gender != "" {
		db = db.Where("gender = ?", info.Gender)
	}

	// 会员状态搜索（需要子查询）
	if info.HasMembership != "" {
		if info.HasMembership == "true" {
			db = db.Where("EXISTS (SELECT 1 FROM user_memberships WHERE user_memberships.user_id = users.id AND user_memberships.status = 'active')")
		} else if info.HasMembership == "false" {
			db = db.Where("NOT EXISTS (SELECT 1 FROM user_memberships WHERE user_memberships.user_id = users.id AND user_memberships.status = 'active')")
		}
	}

	// 注册时间范围搜索
	if info.StartDate != "" {
		db = db.Where("created_at >= ?", info.StartDate)
	}
	if info.EndDate != "" {
		db = db.Where("created_at <= ?", info.EndDate+" 23:59:59")
	}

	// 关键字搜索（搜索用户名、邮箱、昵称）
	if info.Keyword != "" {
		keyword := "%" + info.Keyword + "%"
		db = db.Where("username LIKE ? OR email LIKE ? OR nickname LIKE ?",
			keyword, keyword, keyword)
	}

	return db
}

// buildOrderConditions 构建排序条件
func (u *UserService) buildOrderConditions(order string, desc bool) string {
	defaultOrder := "created_at DESC"
	if order == "" {
		return defaultOrder
	}
	// 验证排序字段安全性
	allowedOrderFields := map[string]bool{
		"id":             true,
		"username":       true,
		"email":          true,
		"created_at":     true,
		"updated_at":     true,
		"last_login_at":  true,
		"login_count":    true,
		"account_status": true,
	}
	if !allowedOrderFields[order] {
		return defaultOrder
	}
	orderStr := order
	if desc {
		orderStr += " DESC"
	} else {
		orderStr += " ASC"
	}

	// 添加二级排序
	if order != "created_at" {
		orderStr += ", created_at DESC"
	}

	return orderStr
}

func (u *UserService) CountUserBy(conditions ...func(*gorm.DB) *gorm.DB) (count int64, err error) {
	query := global.GVA_DB.Model(&project.User{})
	// 应用所有条件
	for _, condition := range conditions {
		query = condition(query)
	}
	err = query.Count(&count).Error
	return count, err
}

// GetSimpleUser  获取简单用户
func (u *UserService) GetSimpleUser(conditions ...func(*gorm.DB) *gorm.DB) (user project.User, err error) {
	query := global.GVA_DB.Model(&project.User{})
	// 应用所有条件
	for _, condition := range conditions {
		query = condition(query)
	}
	err = query.First(&user).Error
	return user, err
}

// RegisterUser 用户注册
func (u *UserService) RegisterUser(req request.BaseRegisterRequest, clientIP string) error {
	// 邀请码
	var referrerID uint
	referrerID = 0
	if req.InviteCode != "" {
		var referrerUser project.User
		err := global.GVA_DB.Model(&project.User{}).First(&referrerUser).Error
		if err == nil && referrerUser.ID > 0 {
			referrerID = referrerUser.ID
		}
	}
	passwordHash := utils.BcryptHash(req.Password)
	// 处理日期
	ReferralCode, err := u.generateReferralCode()
	if err != nil {
		return errors.New("生成邀请码失败")
	}
	// 创建用户
	manager := utils.NewUsernameGeneratorManager()
	// 随机生成任意风格
	username, _ := manager.GenerateRandom()
	user := project.User{
		UUID:          uuid.New(),
		Username:      username,
		Email:         req.Email,
		Phone:         &req.Phone,
		PasswordHash:  passwordHash,
		AccountStatus: constants.AccountStatusNormal,
		EmailVerified: false,
		PhoneVerified: false,
		ReferrerID:    &referrerID,
		ReferralCode:  &ReferralCode,
		RegisterIP:    &clientIP,
	}

	// 开启事务
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 1. 创建用户
		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		// 2. 创建用户统计记录
		statistics := project.UserStatistics{
			UserID: user.ID,
		}
		if err := tx.Create(&statistics).Error; err != nil {
			return err
		}

		// 3. ✅ 创建团队统计记录
		teamStats := project.TeamStatistics{
			UserID: int64(user.ID),
		}
		if err := tx.Create(&teamStats).Error; err != nil {
			return err
		}
		// 4. ✅ 创建佣金账户记录（新增）
		commissionAccount := project.UserCommissionAccount{
			UserID: user.ID,
		}
		if err := tx.Create(&commissionAccount).Error; err != nil {
			return err
		}
		// 5. 如果有推荐人，更新推荐人的统计
		if referrerID > 0 {
			// 确保推荐人的统计记录存在
			var count int64

			// 检查 user_statistics
			tx.Model(&project.UserStatistics{}).Where("user_id = ?", referrerID).Count(&count)
			if count == 0 {
				refStats := project.UserStatistics{UserID: referrerID}
				tx.Create(&refStats)
			}

			// 检查 team_statistics
			tx.Model(&project.TeamStatistics{}).Where("user_id = ?", referrerID).Count(&count)
			if count == 0 {
				refTeamStats := project.TeamStatistics{UserID: int64(referrerID)}
				tx.Create(&refTeamStats)
			}

			// ✅ 更新 user_statistics
			if err := tx.Model(&project.UserStatistics{}).
				Where("user_id = ?", referrerID).
				UpdateColumn("successful_referrals", gorm.Expr("successful_referrals + 1")).
				Error; err != nil {
				return err
			}

			// ✅ 更新 team_statistics
			if err := tx.Model(&project.TeamStatistics{}).
				Where("user_id = ?", referrerID).
				Updates(map[string]interface{}{
					"total_members": gorm.Expr("total_members + 1"),
					"today_new":     gorm.Expr("today_new + 1"),
				}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// BatchUpdateUserStatus 批量更新用户状态
func (u *UserService) BatchUpdateUserStatus(ids []uint, status string) error {
	return global.GVA_DB.Model(&project.User{}).
		Where("id IN ?", ids).
		Update("account_status", status).Error
}

// ResetUserPassword 重置用户密码
func (u *UserService) ResetUserPassword(id uint) (string, error) {
	// 生成新密码
	newPassword := u.generateRandomPassword()
	passwordHash := utils.BcryptHash(newPassword)
	// 更新密码
	err := global.GVA_DB.Model(&project.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"password_hash": passwordHash,
		}).Error
	return newPassword, err
}

// ChangeUserPassword 用户修改密码
func (u *UserService) ChangeUserPassword(id uint, req request.ChangeUserPasswordRequest) error {
	var user project.User
	err := global.GVA_DB.Model(&project.User{}).Where("id = ? ", id).First(&user).Error
	if err != nil {
		global.GVA_LOG.Error("用户修改密码，获取登录用户信息失败", zap.Error(err), zap.Any("userId", id))
		return errors.New("获取登录用户信息失败")
	}
	if !utils.BcryptCheck(req.OldPassword, user.PasswordHash) {
		return errors.New("旧密码不正确")
	}
	// 生成新密码
	passwordHash := utils.BcryptHash(req.NewPassword)
	// 更新密码
	err = global.GVA_DB.Model(&project.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"password_hash": passwordHash,
		}).Error
	return err
}

// ApplyWithdraw 用户提现申请
// ApplyWithdraw 申请提现
func (u *UserService) ApplyWithdraw(userID uint, req request.UserWithdrawRequest) error {
	// 1. 验证请求参数
	if err := req.Validate(); err != nil {
		return err
	}

	// 2. 获取提现规则配置
	config, err := systemConfigService.GetConfig("commission")
	if err != nil {
		return errors.New("获取提现配置失败")
	}

	// 解析配置
	withdrawConfig, err := utils.ParseWithdrawConfig(config)
	if err != nil {
		return errors.New("提现配置格式错误")
	}

	// 3. 验证提现金额范围
	if req.Amount < withdrawConfig.MinWithdraw {
		return fmt.Errorf("提现金额不能低于 %.2f 元", withdrawConfig.MinWithdraw)
	}
	if req.Amount > withdrawConfig.MaxWithdraw {
		return fmt.Errorf("提现金额不能超过 %.2f 元", withdrawConfig.MaxWithdraw)
	}

	// 4. 验证提现方式是否支持
	if !utils.Contains(withdrawConfig.WithdrawMethods, req.WithdrawType) {
		return errors.New("不支持该提现方式")
	}

	// 5. 使用事务处理提现流程
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 5.1 检查今日提现次数
		today := time.Now().Truncate(24 * time.Hour)
		var todayCount int64
		if err := tx.Model(&project.WithdrawRecord{}).
			Where("user_id = ? AND create_time >= ?", userID, today).
			Count(&todayCount).Error; err != nil {
			return errors.New("查询提现记录失败")
		}

		if todayCount >= int64(withdrawConfig.DailyWithdrawCount) {
			return fmt.Errorf("今日提现次数已达上限(%d次)", withdrawConfig.DailyWithdrawCount)
		}

		// 5.2 查询用户佣金账户（加锁）
		var account project.UserCommissionAccount
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).
			First(&account).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("佣金账户不存在")
			}
			return errors.New("查询账户信息失败")
		}

		// 5.3 检查可提现余额
		if account.AvailableAmount < req.Amount {
			return fmt.Errorf("可提现余额不足，当前余额：%.2f 元", account.AvailableAmount)
		}

		// 5.4 计算手续费和实际到账金额
		fee := req.Amount * withdrawConfig.WithdrawFee / 100
		actualAmount := req.Amount - fee

		// 5.5 冻结提现金额（从可用金额转到冻结金额）
		if err := tx.Model(&account).Updates(map[string]interface{}{
			"available_amount": gorm.Expr("available_amount - ?", req.Amount),
			"frozen_amount":    gorm.Expr("frozen_amount + ?", req.Amount),
			"updated_at":       time.Now(),
		}).Error; err != nil {
			return errors.New("冻结提现金额失败")
		}
		// 5.6 生成提现单号
		withdrawNo := utils.GenerateWithdrawNo(userID)
		// 5.7 获取账户信息
		accountName, accountNo := req.GetAccountInfo()
		// 5.8 创建提现记录
		record := project.WithdrawRecord{
			UserID:       int64(userID),
			WithdrawNo:   withdrawNo,
			Amount:       req.Amount,
			Fee:          fee,
			ActualAmount: actualAmount,
			WithdrawType: req.WithdrawType,
			AccountName:  &accountName,
			AccountNo:    &accountNo,
			Status:       project.WithdrawStatusPending,
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		}
		if err := tx.Create(&record).Error; err != nil {
			return errors.New("创建提现记录失败")
		}
		return nil
	})
}

// GetWithdrawRecord 获取提现记录（带筛选）
func (u *UserService) GetWithdrawRecord(userID uint, req request.WithdrawRecordRequest) (response.WithdrawRecordListResp, error) {
	var result response.WithdrawRecordListResp
	var records []project.WithdrawRecord

	// 2. 计算分页参数
	limit := req.PageSize
	offset := req.PageSize * (req.Page - 1)

	// 3. 构建基础查询
	db := global.GVA_DB.Model(&project.WithdrawRecord{}).
		Where("user_id = ?", userID)

	// 4. 添加状态筛选
	if statusConditions := req.GetStatusCondition(); statusConditions != nil {
		db = db.Where("status IN ?", statusConditions)
	}

	// 5. 添加时间筛选
	if startTime, endTime := req.GetTimeRange(); startTime != nil {
		if endTime != nil {
			db = db.Where("create_time BETWEEN ? AND ?", startTime, endTime)
		} else {
			db = db.Where("create_time >= ?", startTime)
		}
	}

	// 6. 查询总数
	if err := db.Count(&result.Total).Error; err != nil {
		return result, err
	}

	// 7. 如果没有数据，直接返回
	if result.Total == 0 {
		result.List = []response.WithdrawRecordResp{}
		return result, nil
	}

	// 8. 查询列表数据
	if err := db.Order("create_time DESC").
		Limit(limit).
		Offset(offset).
		Find(&records).Error; err != nil {
		return result, err
	}

	// 9. 转换为响应格式
	result.List = make([]response.WithdrawRecordResp, 0, len(records))
	for _, record := range records {
		resp := response.WithdrawRecordResp{
			ID:           record.ID,
			WithdrawNo:   record.WithdrawNo,
			Amount:       record.Amount,
			Fee:          record.Fee,
			ActualAmount: record.ActualAmount,
			WithdrawType: record.WithdrawType,
			AccountName:  record.AccountName,
			AccountNo:    record.AccountNo,
			Status:       record.Status,
			RejectReason: record.RejectReason,
			AuditTime:    record.AuditTime,
			CompleteTime: record.CompleteTime,
			Remark:       record.Remark,
			CreateTime:   record.CreateTime,
			UpdateTime:   record.UpdateTime,
		}

		// 账号脱敏
		resp.MaskAccountNo()

		result.List = append(result.List, resp)
	}

	// 10. 查询统计数据（累计提现金额和次数）
	var stats struct {
		TotalAmount float64
		TotalCount  int64
	}

	// 统计已完成的提现记录
	statsDB := global.GVA_DB.Model(&project.WithdrawRecord{}).
		Where("user_id = ?", userID).
		Where("status = ?", project.WithdrawStatusCompleted)

	if err := statsDB.Select("COALESCE(SUM(amount), 0) as total_amount, COUNT(*) as total_count").
		Scan(&stats).Error; err != nil {
		// 统计失败不影响列表返回，只记录日志
		global.GVA_LOG.Error("查询提现统计数据失败: " + err.Error())
	}
	result.TotalWithdrawn = stats.TotalAmount
	result.TotalCount = stats.TotalCount
	return result, nil
}

// generateReferralCode 生成推荐码
func (u *UserService) generateReferralCode() (string, error) {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 8
	result := make([]byte, length)
	randomBytes := make([]byte, length)
	// 使用密码学安全的随机数生成器
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	for i := 0; i < length; i++ {
		result[i] = chars[randomBytes[i]%byte(len(chars))]
	}
	return string(result), nil
}

// generateUniqueReferralCode 生成唯一推荐码（检查数据库重复）
func (u *UserService) generateUniqueReferralCode() (string, error) {
	maxAttempts := 10

	for attempts := 0; attempts < maxAttempts; attempts++ {
		code, err := u.generateReferralCode()
		if err != nil {
			return "", err
		}
		// 检查数据库中是否已存在
		var count int64
		if err := global.GVA_DB.Model(&project.User{}).
			Where("referral_code = ?", code).
			Count(&count).Error; err != nil {
			return "", err
		}

		if count == 0 {
			return code, nil
		}
	}
	return "", errors.New("无法生成唯一推荐码，请重试")
}

// GetUserMemberships 获取用户会员记录
func (u *UserService) GetUserMemberships(userID uint) ([]project.UserMembership, error) {
	var memberships []project.UserMembership
	err := global.GVA_DB.Where("user_id = ?", userID).
		Preload("Plan").
		Order("created_at DESC").
		Find(&memberships).Error

	return memberships, err
}

// GetUserOrders 获取用户订单记录
func (u *UserService) GetUserOrders(userID uint) ([]project.Order, error) {
	var orders []project.Order
	err := global.GVA_DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

// 工具方法

// generateRandomPassword 生成随机密码
func (u *UserService) generateRandomPassword() string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, 8)
	rand.Read(bytes)

	for i, b := range bytes {
		bytes[i] = chars[b%byte(len(chars))]
	}

	return string(bytes)
}

// Login 用户登录
func (u *UserService) Login(us *project.User, clientIP, userAgent string) (userInter *project.User, err error) {
	if nil == global.GVA_DB {
		return nil, fmt.Errorf("db not init")
	}
	var user project.User
	err = global.GVA_DB.Where("phone = ?", us.Phone).First(&user).Error
	// 用户不存在
	if err != nil {
		return nil, errors.New("用户名不存在或密码错误")
	}
	// ==================== 第一步：检查账户状态 ====================
	if err := u.checkAccountStatus(&user); err != nil {
		return nil, err
	}

	// ==================== 第二步：验证密码 ====================
	if ok := utils.BcryptCheck(us.PasswordHash, user.PasswordHash); !ok {
		// 密码错误，记录失败
		return nil, u.handleLoginFailure(&user, clientIP)
	}

	// ==================== 第三步：登录成功 ====================
	u.handleLoginSuccess(&user, clientIP, userAgent)
	return &user, err
}

// checkAccountStatus 检查账户状态
func (u *UserService) checkAccountStatus(user *project.User) error {
	// 如果是临时锁定状态，检查是否已过期
	if user.AccountStatus == constants.AccountStatusLocked && user.StatusExpireAt != nil {
		if time.Now().After(*user.StatusExpireAt) {
			// 自动解锁
			user.AccountStatus = constants.AccountStatusNormal
			user.StatusExpireAt = nil
			user.StatusReason = ""
			user.FailedLoginAttempts = 0
			global.GVA_DB.Save(user)
			return nil
		}
		// 还在锁定期内
		remainingMinutes := time.Until(*user.StatusExpireAt).Minutes()
		return fmt.Errorf("账户已被锁定，请在 %.0f 分钟后重试", remainingMinutes)
	}

	// 检查其他阻止登录的状态
	if !user.AccountStatus.CanLogin() {
		reason := user.AccountStatus.GetBlockReason()
		if user.StatusReason != "" {
			reason = reason + ": " + user.StatusReason
		}
		return errors.New(reason)
	}

	return nil
}

// handleLoginFailure 处理登录失败
func (u *UserService) handleLoginFailure(user *project.User, clientIP string) error {
	// 增加失败次数
	user.FailedLoginAttempts++
	now := time.Now()
	user.LastFailedLoginAt = &now
	user.LastFailedLoginIP = clientIP

	// 检查是否需要锁定账户
	maxAttempts := 10 // 可以放到配置文件
	if user.FailedLoginAttempts >= uint(maxAttempts) {
		// 锁定30分钟
		lockDuration := 30 * time.Minute
		lockUntil := time.Now().Add(lockDuration)

		user.AccountStatus = constants.AccountStatusLocked
		user.StatusReason = "连续登录失败次数过多"
		user.StatusExpireAt = &lockUntil

		global.GVA_DB.Save(user)

		// 发送安全警告通知（异步）
		go u.sendSecurityAlert(user, "您的账户因多次登录失败已被临时锁定30分钟")

		global.GVA_LOG.Warn("账户已被锁定",
			zap.String("phone", *user.Phone),
			zap.Uint("failed_attempts", user.FailedLoginAttempts),
			zap.String("ip", clientIP),
		)

		return errors.New("账户已被锁定30分钟，如非本人操作请及时修改密码")
	}

	// 保存失败记录
	global.GVA_DB.Save(user)

	global.GVA_LOG.Warn("登录失败",
		zap.String("phone", *user.Phone),
		zap.Uint("failed_attempts", user.FailedLoginAttempts),
		zap.String("ip", clientIP),
	)
	frequency := uint(maxAttempts) - user.FailedLoginAttempts
	return fmt.Errorf("密码错误，还可尝试 %d 次", frequency)
}

// handleLoginSuccess 处理登录成功
func (u *UserService) handleLoginSuccess(user *project.User, clientIP, userAgent string) {
	// 记录之前的失败次数（用于安全提醒）
	previousFailedAttempts := user.FailedLoginAttempts
	previousFailedLoginAt := user.LastFailedLoginAt // ✅ 保存失败时间
	// 重置失败记录
	user.FailedLoginAttempts = 0
	user.LastFailedLoginAt = nil
	user.LastFailedLoginIP = ""

	// 如果之前是锁定状态，解锁
	if user.AccountStatus == constants.AccountStatusLocked {
		user.AccountStatus = constants.AccountStatusNormal
		user.StatusExpireAt = nil
		user.StatusReason = ""
	}

	// 更新登录成功信息
	now := time.Now()
	user.LastLoginAt = &now
	user.LastLoginIP = &clientIP
	user.LastLoginDevice = &userAgent
	user.LoginCount++
	global.GVA_DB.Save(user)
	// 记录日志
	global.GVA_LOG.Info("用户登录成功",
		zap.String("phone", *user.Phone),
		zap.String("ip", clientIP),
		zap.Uint("login_count", user.LoginCount),
	)
	// 如果之前有多次失败尝试，发送安全提醒
	if previousFailedAttempts > 3 && previousFailedLoginAt != nil {
		go u.sendSecurityAlert(user, fmt.Sprintf("检测到 %d 次失败登录尝试，最后一次在 %s",
			previousFailedAttempts,
			previousFailedLoginAt.Format("2006-01-02 15:04:05"), // ✅ 使用保存的值
		))
	}
}

// sendSecurityAlert 发送安全警告
func (u *UserService) sendSecurityAlert(user *project.User, message string) {
	// TODO: 实现短信/邮件通知
	global.GVA_LOG.Info("发送安全警告",
		zap.String("phone", *user.Phone),
		zap.String("message", message),
	)
}
