package web

import (
	"ApkAdmin/global"
	"ApkAdmin/model/common/response"
	"ApkAdmin/model/project/request"
	"ApkAdmin/service/project"
	"ApkAdmin/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SystemAnnouncementApi struct {
}

// GetSystemAnnouncementList 获取系统公告列表
func (a SystemAnnouncementApi) GetSystemAnnouncementList(c *gin.Context) {
	var req request.ListAnnouncementRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		global.GVA_LOG.Error("获取系统公告列表失败", zap.Error(err))
		response.FailWithMessage("获取系统公告列表失败"+err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := systemAnnouncementService.PageInfoAnnouncement(req)
	if err != nil {
		global.GVA_LOG.Error("获取系统公告列表失败!", zap.Error(err))
		response.FailWithMessage("获取系统公告列表失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, "获取系统公告列表成功", c)

}

// GetSystemAnnouncementDetail 获取系统公告详情
func (a SystemAnnouncementApi) GetSystemAnnouncementDetail(c *gin.Context) {
	var req request.GetById
	err := c.ShouldBindQuery(&req)
	if err != nil {
		global.GVA_LOG.Error("获取系统公告详情请求错误", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	announcement, err := systemAnnouncementService.GetAnnouncement(
		project.WithID(req.Uint()),
		project.WithStatus(1),
	)
	if err != nil {
		global.GVA_LOG.Error("获取系统公告详情失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(announcement, "success", c)
	return
}
