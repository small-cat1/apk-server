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

type AppPackageApi struct {
}

func (a *AppPackageApi) FirstAppPackage(c *gin.Context) {
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
	data, err := AppPackageService.GetAppPackageManual(id)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

func (a *AppPackageApi) GetAppPackageList(c *gin.Context) {
	var pageInfo request2.AppPackageListRequest
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
	countryList, total, err := AppPackageService.GetAppPackageList(pageInfo, pageInfo.OrderKey, pageInfo.Desc)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     countryList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

func (a *AppPackageApi) AddAppPackage(c *gin.Context) {
	var req request2.AppPackageCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("添加应用安装包参数错误!", zap.Error(err))
		response.FailWithMessage("添加应用安装包参数错误，"+err.Error(), c)
		return
	}
	// 验证创建请求
	if err := req.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	userID := utils.GetUserID(c)
	if err := AppPackageService.CreateAppPackage(userID, req); err != nil {
		global.GVA_LOG.Error("创建应用安装包失败!", zap.Error(err))
		response.FailWithMessage("创建应用安装包失败："+err.Error(), c)
		return
	}
	response.OkWithMessage("创建应用安装包成功", c)
}

func (a *AppPackageApi) UpdateAppPackage(c *gin.Context) {
	var req request2.AppPackageUpdateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		global.GVA_LOG.Error("参数错误!", zap.Error(err))
		response.FailWithMessage("参数错误", c)
		return
	}
	// 验证创建请求
	if err := req.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	useID := utils.GetUserID(c)
	err = AppPackageService.UpdateAppPackage(useID, &req)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
func (a *AppPackageApi) DeleteAppPackage(c *gin.Context) {
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
	err = AppPackageService.DeleteAppPackage(info.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

func (a *AppPackageApi) BatchUpdateAppPackageStatus(c *gin.Context) {
	var req request2.ApkBatchUpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("参数错误!", zap.Error(err))
		response.FailWithMessage("参数错误", c)
		return
	}

	if len(req.IDs) == 0 {
		response.FailWithMessage("请选择要更新的安装包", c)
		return
	}
	useID := utils.GetUserID(c)
	if err := AppPackageService.BatchUpdateApkStatus(useID, req.IDs, req.Status); err != nil {
		global.GVA_LOG.Error("批量更新状态失败!", zap.Error(err))
		response.FailWithMessage("批量更新状态失败："+err.Error(), c)
		return
	}

	response.OkWithMessage("批量更新状态成功", c)
}
