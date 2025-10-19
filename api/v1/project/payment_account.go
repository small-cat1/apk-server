package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/common/response"
	"ApkAdmin/model/project/request"
	"ApkAdmin/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PaymentAccountApi struct{}

// CreatePaymentAccount 创建支付账号
// @Tags PaymentAccount
// @Summary 创建支付账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PaymentAccountCreate true "创建支付账号"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /paymentAccounts [post]
func (p *PaymentAccountApi) CreatePaymentAccount(c *gin.Context) {
	var req request.CreatePaymentAccountReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.PaymentAccountVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = paymentAccountService.CreatePaymentAccount(&req)
	if err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdatePaymentAccount 更新支付账号
// @Tags PaymentAccount
// @Summary 更新支付账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PaymentAccountUpdate true "更新支付账号"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /paymentAccounts [put]
func (p *PaymentAccountApi) UpdatePaymentAccount(c *gin.Context) {
	var req request.PaymentAccountUpdate
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = utils.Verify(req, utils.PaymentAccountUpdateVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = paymentAccountService.UpdatePaymentAccount(&req)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeletePaymentAccount 删除支付账号
// @Tags PaymentAccount
// @Summary 删除支付账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PaymentAccountDelete true "删除支付账号"
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /paymentAccounts [delete]
func (p *PaymentAccountApi) DeletePaymentAccount(c *gin.Context) {
	var req request.PaymentAccountDelete
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = paymentAccountService.DeletePaymentAccount(req.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// GetPaymentAccountList 获取支付账号列表
// @Tags PaymentAccount
// @Summary 获取支付账号列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PaymentAccountSearch true "获取支付账号列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "获取成功"
// @Router /paymentAccounts [get]
func (p *PaymentAccountApi) GetPaymentAccountList(c *gin.Context) {
	var pageInfo request.PaymentAccountSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	list, total, err := paymentAccountService.GetPaymentAccountInfoList(pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取成功", c)
}

// GetPaymentAccount 获取单个支付账号详情
// @Tags PaymentAccount
// @Summary 获取单个支付账号详情
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PaymentAccountById true "获取支付账号详情"
// @Success 200 {object} response.Response{data=project.PaymentAccount,msg=string} "获取成功"
// @Router /paymentAccounts/detail [get]
func (p *PaymentAccountApi) GetPaymentAccount(c *gin.Context) {
	var req request.PaymentAccountById
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	account, err := paymentAccountService.GetPaymentAccount(req.ID)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(account, "获取成功", c)
}

// UpdateAccountWeight 更新账号权重
// @Tags PaymentAccount
// @Summary 更新账号权重
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PaymentAccountWeight true "更新账号权重"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /paymentAccounts/weight [put]
func (p *PaymentAccountApi) UpdateAccountWeight(c *gin.Context) {
	var req request.PaymentAccountWeight
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = paymentAccountService.UpdateAccountWeight(req.ID, req.Weight)
	if err != nil {
		global.GVA_LOG.Error("更新权重失败!", zap.Error(err))
		response.FailWithMessage("更新权重失败", c)
		return
	}
	response.OkWithMessage("更新权重成功", c)
}

// ResetDailyAmount 重置日交易金额
// @Tags PaymentAccount
// @Summary 重置日交易金额
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PaymentAccountById true "重置日交易金额"
// @Success 200 {object} response.Response{msg=string} "重置成功"
// @Router /paymentAccounts/reset-daily [put]
func (p *PaymentAccountApi) ResetDailyAmount(c *gin.Context) {
	var req request.PaymentAccountById
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err = paymentAccountService.ResetDailyAmount(req.ID)
	if err != nil {
		global.GVA_LOG.Error("重置日交易金额失败!", zap.Error(err))
		response.FailWithMessage("重置日交易金额失败", c)
		return
	}
	response.OkWithMessage("重置日交易金额成功", c)
}
