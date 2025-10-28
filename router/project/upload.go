package project

import (
	api "ApkAdmin/api/v1"
	"ApkAdmin/middleware"
	"github.com/gin-gonic/gin"
)

type UploadRoute struct {
}

func (u UploadRoute) InitUploadRoute(Router *gin.RouterGroup) {
	uploadApi := api.ApiGroupApp.ProjectApiGroup.UploadApi

	// 需要认证的路由组
	userRouter := Router.Group("upload").Use(middleware.OperationRecord())
	{
		userRouter.POST("getUploadSignature", uploadApi.GetUploadSignature)
	}
}
