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
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"math/rand"
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
		//return &projectRes.DownloadResp{
		//	CanDownload: true,
		//	PackageUrl:  "http://www.baidu.com",
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
func (a AppApi) buildDownloadResp(isVip, canDownload bool, platform constants.Platform, appPackage *projectModel.AppPackage, reason string) *projectRes.DownloadResp {
	resp := &projectRes.DownloadResp{
		CanDownload:    canDownload,
		DownloadReason: reason,
	}
	if !canDownload {
		return resp
	}
	// 根据平台类型设置不同的返回字段
	if platform == constants.PlatformIOS {
		if isVip {
			resp.PackageUrl = *appPackage.FileURL
		} else {
			resp.PackageDetail = a.getFreeIOSAccount()
		}
	} else {
		// Android 返回直接下载链接
		if appPackage.FileURL != nil {
			resp.PackageUrl = *appPackage.FileURL
		}
	}
	return resp
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
