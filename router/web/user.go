package web

import "github.com/gin-gonic/gin"

type UserRouter struct {
}

func (r *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	router := Router.Group("user")
	router.POST("/logout", jwtApi.JsonInBlacklist)            // 用户退出登录
	router.GET("/getUserInfo", userApi.GetUserInfo)           // 获取登录用户的信息
	router.POST("/changePassword", userApi.ChangePassword)    // 修改用户登录密码
	router.POST("/withdraw", userApi.Withdraw)                // 用户提现
	router.GET("/withdraw-config", userApi.GetWithdrawConfig) // 获取用户提现配置

}
