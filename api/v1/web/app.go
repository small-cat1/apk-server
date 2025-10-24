package web

import (
	"ApkAdmin/constants"
	"ApkAdmin/global"
	"ApkAdmin/model/common/response"
	projectModel "ApkAdmin/model/project"
	"ApkAdmin/model/project/request"
	projectRes "ApkAdmin/model/project/response"
	"ApkAdmin/service/project"
	"ApkAdmin/utils"
	"ApkAdmin/utils/upload"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math/rand"
	"strings"
	"time"
)

type AppApi struct {
}

func (a AppApi) ListHotOrRecommendApp(c *gin.Context) {
	appList, err := AppService.GetHotOrRecommendApp()
	if err != nil {
		global.GVA_LOG.Error("获取热门或推荐应用失败!", zap.Error(err))
		response.FailWithMessage("获取热门或推荐应用失败，"+err.Error(), c)
		return
	}
	response.OkWithDetailed(appList, "获取成功", c)
}

func (a AppApi) GetFilterApps(c *gin.Context) {
	var req request.FilterAppRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		global.GVA_LOG.Error("获取分类列表应用失败!", zap.Error(err))
		response.FailWithMessage("获取分类列表应用失败", c)
		return
	}
	err = req.Validate()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := AppService.FilterAppsByCategory(req)
	if err != nil {
		global.GVA_LOG.Error("获取分类列表应用失败!", zap.Error(err))
		response.FailWithMessage("获取分类列表应用失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取分类列表应用成功", c)
}

// GetAccountAppsListByCategory 根据分类获取应用账号
func (a AppApi) GetAccountAppsListByCategory(c *gin.Context) {
	var req request.FilterAccountAppRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		global.GVA_LOG.Error("获取账号应用列表失败!", zap.Error(err))
		response.FailWithMessage("获取账号应用列表失败", c)
		return
	}
	err = req.Validate()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := AppService.GetAppByAccountCategory(req)
	if err != nil {
		global.GVA_LOG.Error("获取账号应用列表失败!", zap.Error(err))
		response.FailWithMessage("获取账号应用列表失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取账号应用列表成功", c)
}

func (a AppApi) SearchApps(c *gin.Context) {
	var req request.SearchAppRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		global.GVA_LOG.Error("搜索应用失败!", zap.Error(err))
		response.FailWithMessage("搜索应用失败", c)
		return
	}
	applicationList, total, err := AppService.SearchApps(req)
	if err != nil {
		global.GVA_LOG.Error("搜索应用失败!", zap.Error(err))
		response.FailWithMessage("搜索应用失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(map[string]interface{}{
		"list":  applicationList,
		"total": total,
	}, "获取成功", c)
}

func (a AppApi) DownloadApp(c *gin.Context) {
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

	// 确定设备平台
	platform := a.detectPlatform(c)
	if platform == "" {
		response.FailWithMessage("不支持的设备类型", c)
		return
	}

	// 处理下载逻辑
	resp, err := a.handleDownloadLogic(c, req.AppId, platform)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

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

// handleDownloadLogic 处理下载逻辑
func (a AppApi) handleDownloadLogic(c *gin.Context, appID uint, platform constants.Platform) (*projectRes.DownloadResp, error) {
	// 获取应用信息
	appInfo, err := AppService.GetAppDetail(appID, platform)
	if err != nil {
		global.GVA_LOG.Error("获取应用安装包失败", zap.Error(err), zap.String("platform", platform.String()))
		return nil, fmt.Errorf("获取应用%s安装包失败", platform.String())
	}

	// 检查是否有可用的安装包
	if len(appInfo.Packages) == 0 {
		return nil, fmt.Errorf("%s设备下暂无支持的安装包", platform.String())
	}

	appPackage := appInfo.Packages[0]

	// 如果是免费应用，直接返回下载信息
	if appInfo.IsFree != nil && *appInfo.IsFree {
		return a.buildDownloadResp(false, true, platform, &appPackage, "success"), nil
	}
	// 收费应用，检查用户权限
	return a.checkUserPermission(c, platform, &appPackage)
}

// checkUserPermission 检查用户下载权限
func (a AppApi) checkUserPermission(c *gin.Context, platform constants.Platform, appPackage *projectModel.AppPackage) (*projectRes.DownloadResp, error) {
	userID := utils.GetUserID(c)
	userDetail, err := UserService.GetUserDetail(project.WithID(userID))
	if err != nil {
		global.GVA_LOG.Error("获取用户信息失败", zap.Error(err))
		return nil, fmt.Errorf("获取用户信息失败")
	}

	// 未付费用户
	if len(userDetail.Memberships) == 0 {
		//url, err := a.GenerateApkDownloadUrl(
		//	"private/package/2025-10-23/Google Chrome_141.0.7390.43_APKPure.apk",
		//	"Google Chrome_141.0.7390.43_APKPure.apk",
		//	300,
		//)
		//if err != nil {
		//	return nil, err
		//}
		//return &projectRes.DownloadResp{
		//	CanDownload: true,
		//	PackageUrl:  url,
		//	//PackageDetail:  "账号georgdowbigginpksu4417@gmail.com密码Aa2501d 密保答案：mm55----mm----mm77----1990年1月1日",
		//	PackageDetail:  "测试test@gmail.com密码Aa123456",
		//	DownloadReason: "success",
		//}, nil
		return &projectRes.DownloadResp{
			CanDownload:    false,
			DownloadReason: "普通用户无法下载，请升级VIP后下载",
		}, nil
	}

	// 检查会员权限
	hasPermission := a.validateUserMembership(userDetail.Memberships, platform)
	if hasPermission {
		//vip用户从用户套餐的详情里面获取账号
		appPackage.FileURL = &userDetail.Memberships[0].Detail
		return a.buildDownloadResp(true, true, platform, appPackage, "success"), nil
	}

	return &projectRes.DownloadResp{
		CanDownload:    false,
		DownloadReason: "用户暂无权限下载,请升级套餐",
	}, nil
}

// validateUserMembership 验证用户会员权限
func (a AppApi) validateUserMembership(memberships []projectModel.UserMembership, platform constants.Platform) bool {
	for _, membership := range memberships {
		// 判断会员是否过期，自动更新状态
		if membership.IsExpired() {
			if err := userMembershipService.UpdateUserMembership(map[string]interface{}{
				"status": constants.MembershipStatusExpired,
			}, project.WithID(membership.ID)); err != nil {
				global.GVA_LOG.Error("更新会员套餐状态失败", zap.Error(err), zap.Any("membership", membership))
			}
			continue
		}

		// 检查会员套餐是否支持指定平台
		if membership.SupportsPlatform(platform.String()) {
			return true
		}
	}
	return false
}

// buildDownloadResp 构建下载响应
func (a *AppApi) buildDownloadResp(
	isVip bool,
	canDownload bool,
	platform constants.Platform,
	appPackage *projectModel.AppPackage,
	reason string,
) *projectRes.DownloadResp {
	resp := &projectRes.DownloadResp{
		CanDownload:    canDownload,
		DownloadReason: reason,
	}

	// 不允许下载时直接返回
	if !canDownload {
		return resp
	}

	// 根据平台处理
	switch platform {
	case constants.PlatformIOS:
		a.handleIOSDownload(resp, isVip, appPackage)
	case constants.PlatformAndroid:
		a.handleAndroidDownload(resp, appPackage)
	}

	return resp
}

// handleIOSDownload 处理iOS下载
func (a *AppApi) handleIOSDownload(resp *projectRes.DownloadResp, isVip bool, appPackage *projectModel.AppPackage) {
	if isVip {
		// VIP用户：生成下载链接
		if url, err := a.getPackageUrl(appPackage); err == nil {
			resp.PackageUrl = url
		}
	} else {
		// 非VIP用户：返回免费账号
		resp.PackageDetail = a.getFreeIOSAccount()
	}
}

// handleAndroidDownload 处理Android下载
func (a *AppApi) handleAndroidDownload(resp *projectRes.DownloadResp, appPackage *projectModel.AppPackage) {
	if url, err := a.getPackageUrl(appPackage); err == nil {
		resp.PackageUrl = url
	}
}

// getPackageUrl 获取安装包下载URL（核心方法）
func (a *AppApi) getPackageUrl(appPackage *projectModel.AppPackage) (string, error) {
	// 1. 空值检查
	if appPackage == nil {
		return "", errors.New("安装包信息为空")
	}

	// 2. 根据OSS类型处理
	switch global.GVA_CONFIG.System.OssType {
	case "aliyun-oss":
		return a.getAliyunOssUrl(appPackage)
	default:
		return a.getFileUrl(appPackage)
	}
}

// getAliyunOssUrl 获取阿里云OSS URL
func (a *AppApi) getAliyunOssUrl(appPackage *projectModel.AppPackage) (string, error) {
	// 检查必要字段
	if appPackage.ObjectName == nil || *appPackage.ObjectName == "" {
		return "", errors.New("OSS对象名称为空")
	}

	objectName := *appPackage.ObjectName

	// 公开文件：直接返回公开URL
	if strings.HasPrefix(objectName, "public/") {
		return a.buildPublicUrl(objectName), nil
	}

	// 私有文件：生成签名URL
	fileName := "package.apk"
	if appPackage.FileName != nil {
		fileName = *appPackage.FileName
	}

	signedUrl, err := a.GenerateApkDownloadUrl(objectName, fileName, 300) // 5分钟
	if err != nil {
		return "", fmt.Errorf("生成签名URL失败: %w", err)
	}

	return signedUrl, nil
}

// getFileUrl 获取文件URL（本地或其他OSS）
func (a *AppApi) getFileUrl(appPackage *projectModel.AppPackage) (string, error) {
	if appPackage.FileURL == nil || *appPackage.FileURL == "" {
		return "", errors.New("文件URL为空")
	}
	return *appPackage.FileURL, nil
}

// buildPublicUrl 构建公开文件URL
func (a *AppApi) buildPublicUrl(objectName string) string {
	return fmt.Sprintf("https://%s.%s/%s",
		global.GVA_CONFIG.AliyunOSS.BucketName,
		global.GVA_CONFIG.AliyunOSS.Endpoint,
		objectName,
	)
}

// 生成APK下载的签名URL
func (a AppApi) GenerateApkDownloadUrl(objectName string, fileName string, expireSeconds int64) (string, error) {
	// 设置强制下载
	options := []oss.Option{
		oss.ResponseContentDisposition(fmt.Sprintf(`attachment; filename="%s"`, fileName)),
	}
	bucket, err := upload.NewBucket()
	if err != nil {
		global.GVA_LOG.Error("functiosn AliyunOSS.NewBucket() Failed", zap.Any("err", err.Error()))
		return "", errors.New("function AliyunOSS.NewBucket() Failed, err:" + err.Error())
	}
	// 生成签名URL（有效期1小时）
	signedUrl, err := bucket.SignURL(objectName, oss.HTTPGet, expireSeconds, options...)
	if err != nil {
		return "", err
	}

	return signedUrl, nil
}

func (a AppApi) getFreeIOSAccount() string {
	// 获取配置
	data, err := websiteConfigService.GetConfigByKey("website", "ios_account")
	if err != nil {
		global.GVA_LOG.Error("获取免费IOS账号失败", zap.Error(err))
		return ""
	}

	// 解析账号列表
	var accounts []string
	accountsByte := data.(string)
	err = json.Unmarshal([]byte(accountsByte), &accounts)
	if err != nil {
		global.GVA_LOG.Error("解析IOS账号失败", zap.Error(err))
		return ""
	}

	// 检查是否有账号
	if len(accounts) == 0 {
		global.GVA_LOG.Warn("IOS账号列表为空")
		return ""
	}

	// 过滤掉空账号
	var validAccounts []string
	for _, account := range accounts {
		if account != "" {
			validAccounts = append(validAccounts, account)
		}
	}
	if len(validAccounts) == 0 {
		global.GVA_LOG.Warn("没有有效的IOS账号")
		return ""
	}
	// 随机选取一个账号
	rand.NewSource(time.Now().UnixNano())
	randomIndex := rand.Intn(len(validAccounts))
	selectedAccount := validAccounts[randomIndex]
	return selectedAccount
}
