package web

import (
	"ApkAdmin/middleware"
	"github.com/gin-gonic/gin"
)

type BaseRouter struct {
}

func (r BaseRouter) InitBaseRouter(Router *gin.RouterGroup) {
	Router.GET("customer-service/config", baseApi.CustomerServiceConfig) // 获取站点客服配置
	Router.GET("captcha", middleware.CaptchaLimit(), baseApi.Captcha)    // 获取验证码
	//Router.POST("base/register", middleware.RegisterLimit(), baseApi.Register) // 用户注册
	Router.POST("base/register", baseApi.Register) // 用户注册
	Router.POST("base/login", baseApi.Login)       // 用户登录
}
