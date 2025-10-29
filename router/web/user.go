package web

import "github.com/gin-gonic/gin"

type UserRouter struct {
}

func (r *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	router := Router.Group("user")
	router.POST("/logout", jwtApi.JsonInBlacklist)         // 用户退出登录
	router.GET("/getUserInfo", userApi.GetUserInfo)        // 获取登录用户的信息
	router.POST("/changePassword", userApi.ChangePassword) // 修改用户登录密码
}
