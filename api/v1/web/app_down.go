package web

import (
	"ApkAdmin/constants"
	"ApkAdmin/global"
	"ApkAdmin/model/common/response"
	projectModel "ApkAdmin/model/project"
	"ApkAdmin/model/project/request"
	projectRes "ApkAdmin/model/project/response"
	"ApkAdmin/utils"
	"ApkAdmin/utils/upload"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	"strings"
	"time"
)

// DownloadApp 下载应用（优化版）
func (a AppApi) DownloadApp(c *gin.Context) {
	// 1. 参数验证
	var req request.DownloadAppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("下载应用失败，参数错误!", zap.Error(err))
		response.FailWithMessage("下载应用失败，参数错误："+err.Error(), c)
		return
	}

	if err := req.Validate(); err != nil {
		response.FailWithMessage("下载应用失败，参数错误："+err.Error(), c)
		return
	}

	// 2. 检测设备平台
	platform := a.detectPlatform(c)
	if platform == "" {
		response.FailWithMessage("不支持的设备类型", c)
		return
	}

	// 3. 处理下载逻辑
	resp, err := a.handleDownloadLogic(c, req.AppId, platform)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 4. ✅ 记录下载日志（异步）
	go a.recordDownloadLog(c, req.AppId, platform, resp.CanDownload)

	response.OkWithData(resp, c)
}

// detectPlatform 检测设备平台
func (a AppApi) detectPlatform(c *gin.Context) constants.Platform {
	if utils.IsIOS(c) {
		return constants.PlatformIOS
	}
	if utils.IsAndroid(c) {
		return constants.PlatformAndroid
	}
	return ""
}

// handleDownloadLogic 处理下载逻辑（优化版）
func (a AppApi) handleDownloadLogic(c *gin.Context, appID uint, platform constants.Platform) (*projectRes.DownloadResp, error) {
	// 1. 获取应用信息
	appInfo, err := AppService.GetAppDetail(appID, platform)
	if err != nil {
		global.GVA_LOG.Error("获取应用安装包失败",
			zap.Error(err),
			zap.String("platform", platform.String()))
		return nil, fmt.Errorf("获取应用%s安装包失败", platform.String())
	}

	// 2. 检查是否有安装包
	if platform == constants.PlatformAndroid && len(appInfo.Packages) == 0 {
		return nil, fmt.Errorf("%s设备下暂无支持的安装包", platform.String())
	}

	// 3. 免费应用处理
	if appInfo.IsFree != nil && *appInfo.IsFree {
		return a.handleFreeAppDownload(platform, &appInfo)
	}

	// 4. 收费应用处理
	return a.handlePaidAppDownload(c, platform, &appInfo)
}

// ✅ handleFreeAppDownload 处理免费应用下载
func (a AppApi) handleFreeAppDownload(platform constants.Platform, appInfo *projectModel.Application) (*projectRes.DownloadResp, error) {
	switch platform {
	case constants.PlatformIOS:
		account := a.getFreeIOSAccount()
		if account == "" {
			return nil, errors.New("暂无可用的免费iOS账号")
		}
		return &projectRes.DownloadResp{
			CanDownload:    true,
			DownloadReason: "success",
			PackageDetail:  account,
		}, nil

	case constants.PlatformAndroid:
		return a.handleAndroidDownload(&appInfo.Packages[0])

	default:
		return nil, errors.New("不支持的平台")
	}
}

// ✅ handlePaidAppDownload 处理收费应用下载（优化版）
func (a AppApi) handlePaidAppDownload(c *gin.Context, platform constants.Platform, appInfo *projectModel.Application) (*projectRes.DownloadResp, error) {
	userID := utils.GetUserID(c)

	// ✅ 使用轻量级查询，只获取会员信息
	membership, err := a.getUserValidMembership(userID, platform)
	if err != nil {
		// 区分不同的错误类型
		if err.Error() == "no_membership" {
			return &projectRes.DownloadResp{
				CanDownload:    false,
				DownloadReason: "普通用户无法下载，请升级VIP后下载",
			}, nil
		}
		if err.Error() == "platform_not_supported" {
			return &projectRes.DownloadResp{
				CanDownload:    false,
				DownloadReason: "当前会员套餐不支持该平台，请升级套餐",
			}, nil
		}
		if err.Error() == "download_limit_exceeded" {
			return &projectRes.DownloadResp{
				CanDownload:    false,
				DownloadReason: "今日下载次数已用完，请明天再试",
			}, nil
		}
		return nil, err
	}

	// ✅ 根据平台返回下载信息
	switch platform {
	case constants.PlatformIOS:
		return &projectRes.DownloadResp{
			CanDownload:    true,
			DownloadReason: "success",
			PackageDetail:  membership.Detail,
		}, nil

	case constants.PlatformAndroid:
		appPackage := &appInfo.Packages[0]

		// ✅ 增加下载计数（异步）
		go a.incrementDownloadCount(membership.ID)

		return a.handleAndroidDownload(appPackage)

	default:
		return nil, errors.New("不支持的平台")
	}
}

// ✅ getUserValidMembership 获取用户有效会员（优化版）
func (a AppApi) getUserValidMembership(userID uint, platform constants.Platform) (*projectModel.UserMembership, error) {
	var memberships []projectModel.UserMembership

	now := time.Now()

	// ✅ 只查询必要的数据
	err := global.GVA_DB.
		Where("user_id = ?", userID).
		Where("status = ?", constants.MembershipStatusActive).
		Where("(end_date IS NULL OR end_date > ?)", now).
		Preload("Plan").        // 只加载套餐信息
		Order("end_date DESC"). // 优先返回有效期最长的
		Find(&memberships).Error

	if err != nil {
		global.GVA_LOG.Error("查询用户会员失败", zap.Error(err))
		return nil, errors.New("查询会员信息失败")
	}

	// 没有会员
	if len(memberships) == 0 {
		return nil, errors.New("no_membership")
	}

	// ✅ 找到支持该平台的最佳会员
	for _, membership := range memberships {
		// 检查是否支持该平台
		if !membership.SupportsPlatform(platform.String()) {
			continue
		}

		// ✅ 检查下载次数限制
		if !membership.CanDownload(true, true) {
			continue
		}

		// 找到可用的会员
		return &membership, nil
	}

	// 有会员但都不支持该平台
	for _, membership := range memberships {
		if !membership.SupportsPlatform(platform.String()) {
			return nil, errors.New("该会员不支持此平台")
		}
	}

	// 有会员但下载次数用完
	return nil, errors.New("下载次数已经用完")
}

// ✅ incrementDownloadCount 增加下载计数（异步）
func (a AppApi) incrementDownloadCount(membershipID uint) {
	var membership projectModel.UserMembership

	// 加锁查询
	err := global.GVA_DB.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ?", membershipID).
		First(&membership).Error

	if err != nil {
		global.GVA_LOG.Error("查询会员失败", zap.Error(err))
		return
	}

	// 增加计数
	membership.IncrementDownloadCount()

	// 更新数据库
	if err := global.GVA_DB.Save(&membership).Error; err != nil {
		global.GVA_LOG.Error("更新下载计数失败", zap.Error(err))
	}
}

// ✅ recordDownloadLog 记录下载日志（异步）
func (a AppApi) recordDownloadLog(c *gin.Context, appID uint, platform constants.Platform, success bool) {
	UserAgent := c.Request.UserAgent()
	log := projectModel.DownloadLog{
		UserID:    utils.GetUserID(c),
		AppID:     appID,
		Platform:  platform,
		Success:   success,
		IP:        c.ClientIP(),
		UserAgent: &UserAgent,
		CreatedAt: time.Now(),
	}
	if err := global.GVA_DB.Create(&log).Error; err != nil {
		global.GVA_LOG.Error("记录下载日志失败", zap.Error(err))
	}
}

// handleAndroidDownload 处理Android下载
func (a *AppApi) handleAndroidDownload(appPackage *projectModel.AppPackage) (*projectRes.DownloadResp, error) {
	url, err := a.getPackageUrl(appPackage)
	if err != nil {
		global.GVA_LOG.Error("生成安卓下载地址失败", zap.Error(err))
		return &projectRes.DownloadResp{
			CanDownload:    false,
			DownloadReason: "下载地址获取失败",
		}, nil
	}

	return &projectRes.DownloadResp{
		CanDownload:    true,
		DownloadReason: "success",
		PackageUrl:     url,
	}, nil
}

// getPackageUrl 获取安装包下载URL
func (a *AppApi) getPackageUrl(appPackage *projectModel.AppPackage) (string, error) {
	if appPackage == nil {
		return "", errors.New("安装包信息为空")
	}

	switch global.GVA_CONFIG.System.OssType {
	case "aliyun-oss":
		return a.getAliyunOssUrl(appPackage)
	default:
		return a.getFileUrl(appPackage)
	}
}

// getAliyunOssUrl 获取阿里云OSS URL
func (a *AppApi) getAliyunOssUrl(appPackage *projectModel.AppPackage) (string, error) {
	if appPackage.ObjectName == nil || *appPackage.ObjectName == "" {
		return "", errors.New("OSS对象名称为空")
	}

	objectName := *appPackage.ObjectName

	// 公开文件：直接返回公开URL
	if strings.HasPrefix(objectName, "public/") {
		return utils.BuildPublicUrl(objectName), nil
	}

	// 私有文件：生成签名URL
	fileName := "package.apk"
	if appPackage.FileName != nil {
		fileName = *appPackage.FileName
	}

	signedUrl, err := a.GenerateApkDownloadUrl(objectName, fileName, 300)
	if err != nil {
		return "", fmt.Errorf("生成签名URL失败: %w", err)
	}

	return signedUrl, nil
}

// getFileUrl 获取文件URL
func (a *AppApi) getFileUrl(appPackage *projectModel.AppPackage) (string, error) {
	if appPackage.FileURL == nil || *appPackage.FileURL == "" {
		return "", errors.New("文件URL为空")
	}
	return *appPackage.FileURL, nil
}

// GenerateApkDownloadUrl 生成APK下载的签名URL
func (a AppApi) GenerateApkDownloadUrl(objectName string, fileName string, expireSeconds int64) (string, error) {
	options := []oss.Option{
		oss.ResponseContentDisposition(fmt.Sprintf(`attachment; filename="%s"`, fileName)),
	}

	bucket, err := upload.NewBucket()
	if err != nil {
		global.GVA_LOG.Error("NewBucket失败", zap.Error(err))
		return "", errors.New("初始化OSS失败: " + err.Error())
	}

	signedUrl, err := bucket.SignURL(objectName, oss.HTTPGet, expireSeconds, options...)
	if err != nil {
		return "", err
	}

	return signedUrl, nil
}

// getFreeIOSAccount 获取免费iOS账号
func (a AppApi) getFreeIOSAccount() string {
	data, err := websiteConfigService.GetConfigByKey("website", "ios_account")
	if err != nil {
		global.GVA_LOG.Error("获取免费IOS账号失败", zap.Error(err))
		return ""
	}

	accounts, err := utils.ParseAccounts(data)
	if err != nil {
		global.GVA_LOG.Error("解析IOS账号失败",
			zap.Error(err),
			zap.String("dataType", fmt.Sprintf("%T", data)))
		return ""
	}

	validAccounts := utils.FilterValidAccounts(accounts)
	if len(validAccounts) == 0 {
		global.GVA_LOG.Warn("没有有效的IOS账号")
		return ""
	}

	return utils.RandomSelectAccount(validAccounts)
}
