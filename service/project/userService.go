package project

import (
	"ApkAdmin/constants"
	"ApkAdmin/global"
	"ApkAdmin/model/project"
	"ApkAdmin/model/project/request"
	"ApkAdmin/utils"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
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

// GetUserDetail  获取用户详情
func (u *UserService) GetUserDetail(conditions ...func(*gorm.DB) *gorm.DB) (user project.User, err error) {
	query := global.GVA_DB.Model(&project.User{})
	// 应用所有条件
	for _, condition := range conditions {
		query = condition(query)
	}
	err = query.Preload("Statistics").
		Preload("Memberships", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("Referrer").
		First(&user).Error
	return user, err
}

// RegisterUser 用户注册
func (u *UserService) RegisterUser(req request.BaseRegisterRequest, clientIP string) error {
	// 邀请码
	var referrerID uint
	referrerID = 0
	if req.InviteCode != "" {
		var referrerUser project.User
		_ = global.GVA_DB.Model(&project.User{}).First(&referrerUser).Error
		if referrerUser.ID > 0 {
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
		// 创建用户
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		// 创建用户统计记录
		statistics := project.UserStatistics{
			UserID: user.ID,
		}
		if err := tx.Create(&statistics).Error; err != nil {
			return err
		}
		return nil
	})
}

// CreateUser 创建用户
func (u *UserService) CreateUser(req request.CreateUserRequest) error {
	var count int64
	// 检查邮箱是否已存在
	global.GVA_DB.Model(&project.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		return errors.New("邮箱已存在")
	}

	// 检查手机号是否已存在
	if req.Phone != "" {
		global.GVA_DB.Model(&project.User{}).Where("phone = ?", req.Phone).Count(&count)
		if count > 0 {
			return errors.New("手机号已存在")
		}
	}
	passwordHash := utils.BcryptHash(req.Password)
	// 处理日期
	ReferralCode, err := u.generateReferralCode()
	if err != nil {
		return errors.New("生成邀请码失败")
	}
	// 创建用户
	// 创建用户
	manager := utils.NewUsernameGeneratorManager()
	// 随机生成任意风格
	username, _ := manager.GenerateRandom()
	user := project.User{
		UUID:          uuid.New(),
		Username:      username,
		Email:         req.Email,
		PasswordHash:  passwordHash,
		AccountStatus: constants.AccountStatusNormal,
		EmailVerified: req.EmailVerified,
		PhoneVerified: req.PhoneVerified,
		ReferralCode:  &ReferralCode,
	}

	if req.Phone != "" {
		user.Phone = &req.Phone
	}

	// 开启事务
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 创建用户
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		// 创建用户统计记录
		statistics := project.UserStatistics{
			UserID: user.ID,
		}
		if err := tx.Create(&statistics).Error; err != nil {
			return err
		}
		return nil
	})
}

// UpdateUser 更新用户
func (u *UserService) UpdateUser(req request.UpdateUserRequest) error {
	// 检查用户是否存在
	var existingUser project.User
	if err := global.GVA_DB.First(&existingUser, req.ID).Error; err != nil {
		return errors.New("用户不存在")
	}

	var count int64

	// 检查邮箱重复（排除自己）
	global.GVA_DB.Model(&project.User{}).Where("email = ? AND id != ?", req.Email, req.ID).Count(&count)
	if count > 0 {
		return errors.New("邮箱已存在")
	}

	// 检查手机号重复（排除自己）
	if req.Phone != "" {
		global.GVA_DB.Model(&project.User{}).Where("phone = ? AND id != ?", req.Phone, req.ID).Count(&count)
		if count > 0 {
			return errors.New("手机号已存在")
		}
	}
	// 更新数据
	updates := map[string]interface{}{
		"email":          req.Email,
		"account_status": req.AccountStatus,
		"email_verified": req.EmailVerified,
		"phone_verified": req.PhoneVerified,
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	} else {
		updates["phone"] = nil
		updates["phone_verified"] = false
	}

	return global.GVA_DB.Model(&existingUser).Updates(updates).Error
}

// DeleteUser 删除用户
func (u *UserService) DeleteUser(id uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		// 软删除用户
		if err := tx.Delete(&project.User{}, id).Error; err != nil {
			return err
		}
		// 可以选择级联删除相关数据，或者保留用于审计
		// 这里选择保留会员记录和订单记录

		return nil
	})
}

// BatchDeleteUsers 批量删除用户
func (u *UserService) BatchDeleteUsers(ids []uint) error {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&project.User{}, ids).Error
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
func (u *UserService) GetUserOrders(userID uint) ([]project.MembershipOrder, error) {
	var orders []project.MembershipOrder
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
