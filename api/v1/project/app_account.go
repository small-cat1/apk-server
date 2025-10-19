package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/common/response"
	"ApkAdmin/model/project/request"
	"ApkAdmin/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AppAccountApi 应用账号APi
type AppAccountApi struct {
}

func (a AppAccountApi) ListAppAccount(c *gin.Context) {
	var pageInfo request.GetAppAccountListRequest
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
	countryList, total, err := AppAccountService.GetAccountList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取应用账号失败!", zap.Error(err))
		response.FailWithMessage("获取应用账号失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     countryList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取应用账号成功", c)
}

func (a AppAccountApi) GetAppAccountDetail(c *gin.Context) {
	var req request.GetById
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := AppAccountService.GetAccountDetail(uint(req.ID))
	if err != nil {
		global.GVA_LOG.Error("获取应用账号详情失败!", zap.Error(err))
		response.FailWithMessage("获取应用账号详情失败,"+err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}

// ViewAppAccountOrder 查看应用账号订单详情
func (a AppAccountApi) ViewAppAccountOrder(c *gin.Context) {
	var req request.ViewAppAccountOrderRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	data, err := AppAccountService.GetAppAccountOrderDetail(req.AccountId)
	if err != nil {
		global.GVA_LOG.Error("获取应用账号订单详情失败!", zap.Error(err))
		response.FailWithMessage("获取应用账号订单详情失败,"+err.Error(), c)
		return
	}
	response.OkWithData(data, c)
}

func (a AppAccountApi) CreateAppAccount(c *gin.Context) {
	var req request.CreateAppAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("添加应用账号参数错误!", zap.Error(err))
		response.FailWithMessage("添加应用账号参数错误,"+err.Error(), c)
		return
	}
	userId := utils.GetUserID(c)
	err := AppAccountService.CreateAccount(req, userId)
	if err != nil {
		global.GVA_LOG.Error("添加应用账号失败!", zap.Error(err))
		response.FailWithMessage("添加应用账号失败："+err.Error(), c)
		return
	}
	response.OkWithMessage("添加应用账号成功", c)
}

func (a AppAccountApi) UpdateAppAccount(c *gin.Context) {
	var req request.UpdateAppAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("更新应用账号参数错误!", zap.Error(err))
		response.FailWithMessage("更新应用账号参数错误,"+err.Error(), c)
		return
	}
	userID := utils.GetUserID(c)
	err := AppAccountService.UpdateAccount(req, userID)
	if err != nil {
		global.GVA_LOG.Error("获取应用账号详情失败!", zap.Error(err))
		response.FailWithMessage("获取应用账号详情失败,"+err.Error(), c)
		return
	}
	response.OkWithMessage("success", c)
}

func (a AppAccountApi) DeleteAppAccount(c *gin.Context) {

}

func (a AppAccountApi) BatchUpdateStatus(c *gin.Context) {
	var req request.UpdateAccountStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("更新应用账号状态参数错误!", zap.Error(err))
		response.FailWithMessage("更新应用账号状态参数错误,"+err.Error(), c)
		return
	}
	resp, err := AppAccountService.UpdateAccountStatus(req)
	if err != nil {
		global.GVA_LOG.Error("更新应用账号状态失败!", zap.Error(err))
		response.FailWithMessage("更新应用账号状态失败,"+err.Error(), c)
		return
	}
	response.OkWithDetailed(resp, "success", c)
}
