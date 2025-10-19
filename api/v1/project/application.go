package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/common/request"
	"ApkAdmin/model/common/response"
	request2 "ApkAdmin/model/project/request"
	"ApkAdmin/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ApplicationApi struct{}

// CreateApplication 创建应用
func (a *ApplicationApi) CreateApplication(c *gin.Context) {
	var req request2.ApplicationCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("参数错误!", zap.Error(err))
		response.FailWithMessage("参数错误", c)
		return
	}

	// 验证创建请求
	if err := req.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userID := utils.GetUserID(c)
	if err := ApplicationService.CreateApplication(userID, req); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败："+err.Error(), c)
		return
	}

	response.OkWithMessage("创建成功", c)
}

// UpdateApplication 更新应用
func (a *ApplicationApi) UpdateApplication(c *gin.Context) {
	var req request2.ApplicationUpdateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("参数错误!", zap.Error(err))
		response.FailWithMessage("参数错误", c)
		return
	}
	// 验证更新请求
	if err := req.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = ApplicationService.UpdateApplication(&req)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败："+err.Error(), c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeleteApplication 删除应用
func (a *ApplicationApi) DeleteApplication(c *gin.Context) {
	var info request.GetById
	err := c.ShouldBindJSON(&info)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(info, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = ApplicationService.DeleteApplication(info.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败："+err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// BatchDeleteApplications 批量删除应用
func (a *ApplicationApi) BatchDeleteApplications(c *gin.Context) {
	var req request2.ApplicationBatchDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("参数错误!", zap.Error(err))
		response.FailWithMessage("参数错误", c)
		return
	}

	if len(req.IDs) == 0 {
		response.FailWithMessage("请选择要删除的应用", c)
		return
	}

	if err := ApplicationService.BatchDeleteApplications(req.IDs); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败："+err.Error(), c)
		return
	}

	response.OkWithMessage("批量删除成功", c)
}

// BatchUpdateApplicationStatus 批量更新应用状态
func (a *ApplicationApi) BatchUpdateApplicationStatus(c *gin.Context) {
	var req request2.ApplicationBatchUpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("参数错误!", zap.Error(err))
		response.FailWithMessage("参数错误", c)
		return
	}

	if len(req.IDs) == 0 {
		response.FailWithMessage("请选择要更新的应用", c)
		return
	}

	if err := ApplicationService.BatchUpdateApplicationStatus(req.IDs, req.Status); err != nil {
		global.GVA_LOG.Error("批量更新状态失败!", zap.Error(err))
		response.FailWithMessage("批量更新状态失败："+err.Error(), c)
		return
	}

	response.OkWithMessage("批量更新状态成功", c)
}

// CloneApplication 克隆应用
func (a *ApplicationApi) CloneApplication(c *gin.Context) {
	var req request2.ApplicationCloneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("参数错误!", zap.Error(err))
		response.FailWithMessage("参数错误", c)
		return
	}

	// 验证克隆请求
	if err := req.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	uid := utils.GetUserID(c)
	newApp, err := ApplicationService.CloneApplication(uid, req)
	if err != nil {
		global.GVA_LOG.Error("克隆失败!", zap.Error(err))
		response.FailWithMessage("克隆失败："+err.Error(), c)
		return
	}

	response.OkWithDetailed(newApp, "克隆成功", c)
}

// UploadApplicationIcon 上传应用图标
func (a *ApplicationApi) UploadApplicationIcon(c *gin.Context) {
	var req request2.ApplicationUploadIconRequest
	if err := c.ShouldBind(&req); err != nil {
		global.GVA_LOG.Error("参数错误!", zap.Error(err))
		response.FailWithMessage("参数错误", c)
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.FailWithMessage("获取文件失败："+err.Error(), c)
		return
	}

	iconURL, err := ApplicationService.UploadApplicationIcon(req.AppID, file)
	if err != nil {
		global.GVA_LOG.Error("上传图标失败!", zap.Error(err))
		response.FailWithMessage("上传图标失败："+err.Error(), c)
		return
	}

	response.OkWithDetailed(gin.H{"icon_url": iconURL}, "上传成功", c)
}

// GetApplicationList 获取应用列表
func (a *ApplicationApi) GetApplicationList(c *gin.Context) {
	var pageInfo request2.ApplicationListRequest
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	applicationList, total, err := ApplicationService.GetApplicationList(pageInfo, pageInfo.OrderKey, pageInfo.Desc)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	}

	response.OkWithDetailed(response.PageResult{
		List:     applicationList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetApplication 获取应用详情
func (a *ApplicationApi) GetApplication(c *gin.Context) {
	var idInfo request.GetById
	err := c.ShouldBindQuery(&idInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(idInfo, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	id := idInfo.Uint()
	data, err := ApplicationService.GetApplication(id)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败："+err.Error(), c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}
