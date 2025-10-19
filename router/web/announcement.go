package web

import "github.com/gin-gonic/gin"

// AnnouncementRouter 系统公告路由
type AnnouncementRouter struct {
}

func (r *AnnouncementRouter) InitAnnouncementRouter(Router *gin.RouterGroup) {
	Router.GET("getAnnouncementList", systemAnnouncementApi.GetSystemAnnouncementList)
	Router.GET("getAnnouncementDetail", systemAnnouncementApi.GetSystemAnnouncementDetail)
}
