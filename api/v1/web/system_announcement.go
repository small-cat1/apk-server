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
	uid := utils.GetUserID(c)
	list, total, err := systemAnnouncementService.PageInfoAnnouncement(req, uid)
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

// MarkAsRead 标记公告为已读
// @Summary 标记公告为已读
// @Tags 公告
// @Accept json
// @Produce json
// @Param data body request.MarkAnnouncementReadRequest true "公告ID"
// @Success 200 {object} response.Response{msg=string} "标记成功"
// @Router /api/v1/announcement/mark-read [post]
func (a *SystemAnnouncementApi) MarkAsRead(c *gin.Context) {
	var req request.MarkAnnouncementReadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 从 token 中获取用户ID
	userID := utils.GetUserID(c) // 根据你的项目调整获取方式
	if err := systemAnnouncementService.MarkAsRead(int64(userID), req.AnnouncementID); err != nil {
		global.GVA_LOG.Error("用户阅读公告标记失败", zap.Error(err))
		response.FailWithMessage("标记失败: ", c)
		return
	}
	response.OkWithMessage("标记成功", c)
}
