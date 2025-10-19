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

type PaymentProviderApi struct {
}

// GetPaymentProvider 获取单个支付服务商详情
func (a *PaymentProviderApi) GetPaymentProvider(c *gin.Context) {
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
	data, err := PaymentProviderService.GetPaymentProvider(id)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

// GetPaymentProviderList 获取支付服务商列表
func (a *PaymentProviderApi) GetPaymentProviderList(c *gin.Context) {
	var pageInfo request2.PaymentProviderListRequest
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		global.GVA_LOG.Error("获取支付服务商列表失败", zap.Error(err))
		response.FailWithMessage("获取支付服务商失败"+err.Error(), c)
		return
	}
	err = utils.Verify(pageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	providerList, total, err := PaymentProviderService.GetPaymentProviderList(pageInfo, pageInfo.OrderKey, pageInfo.Desc)
	if err != nil {
		global.GVA_LOG.Error("获取支付服务商列表失败!", zap.Error(err))
		response.FailWithMessage("获取支付服务商列表失败"+err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     providerList,
		Total:    total,
		Page:     pageInfo.Page,
		PageSize: pageInfo.PageSize,
	}, "获取支付服务商列表成功", c)
}

// CreatePaymentProvider 创建支付服务商
func (a *PaymentProviderApi) CreatePaymentProvider(c *gin.Context) {
	var req request2.PaymentProviderCreateRequest
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
	if err := PaymentProviderService.CreatePaymentProvider(req); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败："+err.Error(), c)
		return
	}
	response.OkWithMessage("创建成功", c)
}

// UpdatePaymentProvider 更新支付服务商
func (a *PaymentProviderApi) UpdatePaymentProvider(c *gin.Context) {
	var req request2.PaymentProviderUpdateRequest
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
	err = PaymentProviderService.UpdatePaymentProvider(&req)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}

// DeletePaymentProvider 删除支付服务商
func (a *PaymentProviderApi) DeletePaymentProvider(c *gin.Context) {
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
	err = PaymentProviderService.DeletePaymentProvider(uint(info.ID))
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}

// BatchUpdatePaymentProviderStatus 批量更新支付服务商状态
func (a *PaymentProviderApi) BatchUpdatePaymentProviderStatus(c *gin.Context) {
	var req request2.BatchUpdateStatusRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if len(req.IDs) == 0 {
		response.FailWithMessage("请选择要更新的数据", c)
		return
	}
	err = PaymentProviderService.BatchUpdatePaymentProviderStatus(req.IDs, req.Status)
	if err != nil {
		global.GVA_LOG.Error("批量更新状态失败!", zap.Error(err))
		response.FailWithMessage("批量更新状态失败", c)
		return
	}
	response.OkWithMessage("批量更新状态成功", c)
}

// CheckProviderCodeAvailable 检查服务商代码是否可用
func (a *PaymentProviderApi) CheckProviderCodeAvailable(c *gin.Context) {
	var req request2.CheckCodeRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req.Code == "" {
		response.FailWithMessage("服务商代码不能为空", c)
		return
	}
	available, err := PaymentProviderService.CheckProviderCodeAvailable(req.Code, req.ID)
	if err != nil {
		global.GVA_LOG.Error("检查代码可用性失败!", zap.Error(err))
		response.FailWithMessage("检查失败", c)
		return
	}
	response.OkWithDetailed(map[string]bool{"available": available}, "检查成功", c)
}

// UpdateProviderSortOrder 更新服务商排序
func (a *PaymentProviderApi) UpdateProviderSortOrder(c *gin.Context) {
	var req request2.UpdateSortOrderRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = utils.Verify(req, utils.IdVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = PaymentProviderService.UpdateProviderSortOrder(req.ID, req.SortOrder)
	if err != nil {
		global.GVA_LOG.Error("更新排序失败!", zap.Error(err))
		response.FailWithMessage("更新排序失败", c)
		return
	}
	response.OkWithMessage("更新排序成功", c)
}
