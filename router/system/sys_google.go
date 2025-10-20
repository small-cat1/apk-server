package system

import (
	"ApkAdmin/middleware"
	"github.com/gin-gonic/gin"
)

type GoogleAuthRouter struct {
}

func (r GoogleAuthRouter) InitGoogleAuthRouter(Router *gin.RouterGroup) {
	router := Router.Group("google-auth").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("google-auth")
	{
		router.POST("bind", googleAuthApi.BindGoogleAuth)               // 绑定谷歌验证器
		router.POST("verifyGoogleAuth", googleAuthApi.VerifyGoogleAuth) // 验证谷歌验证码
	}
	{
		routerWithoutRecord.GET("qrcode", googleAuthApi.GetGoogleAuthInfo) //获取谷歌验证器
	}
}
