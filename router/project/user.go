package project

import (
	"ApkAdmin/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouter struct {
}

func (r *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	// 需要认证的路由组
	userRouter := Router.Group("user").Use(middleware.OperationRecord())
	userRouterWithoutRecord := Router.Group("user")
	{
		// 管理员接口（需要认证和权限）
		userRouter.POST("createUser", userApi.CreateUser)                       // 创建用户
		userRouter.PUT("updateUser", userApi.UpdateUser)                        // 更新用户
		userRouter.DELETE("removeUser/:id", userApi.DeleteUser)                 // 删除用户
		userRouter.POST("batchDeleteUsers", userApi.BatchDeleteUsers)           // 批量删除用户
		userRouter.POST("batchUpdateUserStatus", userApi.BatchUpdateUserStatus) // 批量更新状态
		userRouter.POST("resetUserPassword", userApi.ResetUserPassword)         // 重置密码
	}
	{
		// 查询接口（需要认证但不记录操作日志）
		userRouterWithoutRecord.GET("getUserList", userApi.GetUserList)                        // 获取用户列表
		userRouterWithoutRecord.GET("getUser/:id", userApi.GetUser)                            // 获取用户详情
		userRouterWithoutRecord.GET("getUserMemberships/:user_id", userApi.GetUserMemberships) // 获取用户会员记录
		userRouterWithoutRecord.GET("getUserOrders/:user_id", userApi.GetUserOrders)           // 获取用户订单记录
	}

}
