package project

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

	list, total, err := systemAnnouncementService.ListAnnouncement(req)
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

// CreateSystemAnnouncement 创建系统公告
func (a SystemAnnouncementApi) CreateSystemAnnouncement(c *gin.Context) {
	var req request.CreateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("创建系统公告参数请求错误", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 进行业务验证
	err := req.Validate()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userID := utils.GetUserID(c)
	err = systemAnnouncementService.CreateAnnouncement(req, userID)
	if err != nil {
		global.GVA_LOG.Error("创建系统公告失败!", zap.Error(err))
		response.FailWithMessage("创建系统公告失败，错误信息："+err.Error(), c)
		return
	}
	response.OkWithMessage("创建系统公告成功", c)
}

// UpdateAnnouncement 更新系统公告
func (a SystemAnnouncementApi) UpdateAnnouncement(c *gin.Context) {
	var req request.UpdateAnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("更新系统公告参数请求错误", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 进行业务验证
	err := req.Validate()
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = systemAnnouncementService.UpdateAnnouncement(req)
	if err != nil {
		global.GVA_LOG.Error("更新系统公告失败!", zap.Error(err))
		response.FailWithMessage("更新系统公告失败，错误信息："+err.Error(), c)
		return
	}
	response.OkWithMessage("更新系统公告成功", c)
}

func (a SystemAnnouncementApi) DeleteAnnouncement(c *gin.Context) {
	var req request.DeleteIDs
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("批量删除系统公告参数请求错误", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 验证
	if err := req.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := systemAnnouncementService.DeleteAnnouncements(req.IDs); err != nil {
		global.GVA_LOG.Error("批量删除系统公告失败", zap.Error(err))
		response.FailWithMessage("批量删除系统公告失败", c)
		return
	}
	response.OkWithMessage("批量删除系统公告成功", c)
	return
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
	announcement, err := systemAnnouncementService.GetAnnouncement(project.WithID(req.Uint()))
	if err != nil {
		global.GVA_LOG.Error("获取系统公告详情失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(announcement, "success", c)
	return
}
