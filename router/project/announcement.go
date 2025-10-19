package project

import (
	"ApkAdmin/middleware"
	"github.com/gin-gonic/gin"
)

type SystemAnnouncementRouter struct {
}

func (s SystemAnnouncementRouter) InitSystemAnnouncementRouter(Router *gin.RouterGroup) {
	router := Router.Group("announcement").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("announcement")
	{
		router.POST("createAnnouncement", systemAnnouncementApi.CreateSystemAnnouncement) // 创建系统公告
		router.PUT("updateAnnouncement", systemAnnouncementApi.UpdateAnnouncement)        // 更新系统公告
		router.DELETE("deleteAnnouncement", systemAnnouncementApi.DeleteAnnouncement)     // 删除系统公告
	}
	{
		routerWithoutRecord.GET("getAnnouncementList", systemAnnouncementApi.GetSystemAnnouncementList) // 系统公告列表
		routerWithoutRecord.GET("findAnnouncement", systemAnnouncementApi.GetSystemAnnouncementDetail)  // 应用公告详情
	}
}
