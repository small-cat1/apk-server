package system

import (
	"ApkAdmin/constants"
	"ApkAdmin/global"
	"ApkAdmin/model/common/response"
	systemReq "ApkAdmin/model/system/request"
	"ApkAdmin/utils"
	"github.com/gin-gonic/gin"
	"github.com/pquerna/otp/totp"
	"go.uber.org/zap"
)

type GoogleAuthApi struct {
}

func (a GoogleAuthApi) GetGoogleAuthInfo(c *gin.Context) {
	uuid := utils.GetUserUuid(c)
	user, err := userService.GetUserInfo(uuid)
	if err != nil {
		global.GVA_LOG.Error("获取用户信息失败", zap.Error(err))
		response.FailWithMessage("获取用户信息失败", c)
		return
	}
	// 检查用户是否已绑定
	if user.GoogleAuthStatus {
		response.OkWithDetailed(gin.H{
			"isBound": true,
		}, "获取成功", c)
		return
	}
	// 生成新的谷歌验证器密钥
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      constants.SystemName,
		AccountName: user.Username,
	})
	if err != nil {
		response.FailWithMessage("生成验证器密钥失败", c)
		return
	}
	response.OkWithDetailed(gin.H{
		"isBound": false,
		"secret":  key.Secret(),
		"qrcode":  key.URL(),
	}, "获取成功", c)
}

// BindGoogleAuth 绑定谷歌验证器
// @Summary 绑定谷歌验证器
// @Router /system/bindGoogleAuth [post]
func (a GoogleAuthApi) BindGoogleAuth(c *gin.Context) {
	var req systemReq.GoogleAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误"+err.Error(), c)
		return
	}
	uuid := utils.GetUserUuid(c)
	// 验证谷歌验证码
	valid := totp.Validate(req.Code, req.Secret)
	if !valid {
		response.FailWithMessage("验证码错误", c)
		return
	}
	user, err := userService.GetUserInfo(uuid)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}
	// 检查用户是否已绑定
	if user.GoogleAuthStatus {
		response.FailWithMessage("用户已绑定，请勿重复操作", c)
		return
	}
	// 绑定谷歌验证器
	err = userService.BindGoogleAuth(uuid, req.Secret)
	if err != nil {
		response.FailWithMessage("绑定失败", c)
		return
	}
	response.OkWithMessage("绑定成功", c)
}

// VerifyGoogleAuth 验证谷歌验证码
// @Summary 验证谷歌验证码
// @Router /system/verifyGoogleAuth [post]
func (a GoogleAuthApi) VerifyGoogleAuth(c *gin.Context) {
	var req systemReq.GoogleAuthVerifyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误"+err.Error(), c)
		return
	}
	// 获取当前用户
	uuid := utils.GetUserUuid(c)
	// 获取用户谷歌验证器密钥
	user, err := userService.GetUserInfo(uuid)
	if err != nil {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}
	if !user.GoogleAuthStatus {
		response.FailWithMessage("未绑定谷歌验证器", c)
		return
	}
	// 验证谷歌验证码
	valid := totp.Validate(req.Code, user.GoogleAuthKey)
	if !valid {
		response.FailWithMessage("验证码错误", c)
		return
	}
	response.OkWithMessage("验证成功", c)
}
