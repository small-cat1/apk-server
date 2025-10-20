package project

import (
	"ApkAdmin/global"
	"ApkAdmin/model/common/response"
	"ApkAdmin/model/project"
	"ApkAdmin/model/project/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CommissionTierApi struct{}

// GetCommissionTierList 获取阶梯等级列表
// @Tags CommissionTier
// @Summary 获取阶梯等级列表
// @Description 获取阶梯等级列表，支持分页和搜索
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param pageSize query int false "每页数量"
// @Param name query string false "等级名称"
// @Param status query int false "状态"
// @Success 200 {object} response.Response{data=response.PageResult} "成功"
// @Router /commissionTier/list [get]
func (a *CommissionTierApi) GetCommissionTierList(c *gin.Context) {
	var req request.GetCommissionTierListRequest
	// 绑定查询参数
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 验证请求参数
	if err := req.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 获取列表
	list, total, err := commissionTierService.GetCommissionTierList(req.Name, req.Status, req.Page, req.PageSize)
	if err != nil {
		global.GVA_LOG.Error("获取阶梯等级列表失败", zap.Error(err))
		response.FailWithMessage("获取列表失败", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:     list,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, c)
}

// GetCommissionTier 获取单个阶梯等级详情
// @Tags CommissionTier
// @Summary 获取阶梯等级详情
// @Description 根据ID获取阶梯等级详情
// @Accept json
// @Produce json
// @Param data body request.GetCommissionTierRequest true "等级ID"
// @Success 200 {object} response.Response{data=project.CommissionTier} "成功"
// @Router /commissionTier/get [post]
func (a *CommissionTierApi) GetCommissionTier(c *gin.Context) {
	var req request.GetCommissionTierRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := req.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	tier, err := commissionTierService.GetCommissionTierById(req.ID)
	if err != nil {
		global.GVA_LOG.Error("获取阶梯等级失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(gin.H{"tier": tier}, c)
}

// CreateCommissionTier 创建阶梯等级
// @Tags CommissionTier
// @Summary 创建阶梯等级
// @Description 创建新的阶梯等级
// @Accept json
// @Produce json
// @Param data body request.CreateCommissionTierRequest true "等级信息"
// @Success 200 {object} response.Response{data=project.CommissionTier} "成功"
// @Router /commissionTier/create [post]
func (a *CommissionTierApi) CreateCommissionTier(c *gin.Context) {
	var req request.CreateCommissionTierRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := req.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 转换为模型
	tier := &project.CommissionTier{
		Name:            req.Name,
		MinSubordinates: req.MinSubordinates,
		Rate:            req.Rate,
		Color:           req.Color,
		Icon:            req.Icon,
		Sort:            req.Sort,
		Status:          &req.Status,
	}

	if err := commissionTierService.CreateCommissionTier(tier); err != nil {
		global.GVA_LOG.Error("创建阶梯等级失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(gin.H{"tier": tier}, c)
}

// UpdateCommissionTier 更新阶梯等级
// @Tags CommissionTier
// @Summary 更新阶梯等级
// @Description 更新阶梯等级信息
// @Accept json
// @Produce json
// @Param data body request.UpdateCommissionTierRequest true "等级信息"
// @Success 200 {object} response.Response "成功"
// @Router /commissionTier/update [post]
func (a *CommissionTierApi) UpdateCommissionTier(c *gin.Context) {
	var req request.UpdateCommissionTierRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := req.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 转换为模型
	tier := &project.CommissionTier{
		ID:              req.ID,
		Name:            req.Name,
		MinSubordinates: req.MinSubordinates,
		Rate:            req.Rate,
		Color:           req.Color,
		Icon:            req.Icon,
		Sort:            req.Sort,
		Status:          &req.Status,
	}

	if err := commissionTierService.UpdateCommissionTier(tier); err != nil {
		global.GVA_LOG.Error("更新阶梯等级失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// DeleteCommissionTiers 删除阶梯等级（支持批量）
// @Tags CommissionTier
// @Summary 删除阶梯等级
// @Description 删除一个或多个阶梯等级
// @Accept json
// @Produce json
// @Param data body request.DeleteCommissionTiersRequest true "等级ID列表"
// @Success 200 {object} response.Response "成功"
// @Router /commissionTier/delete [post]
func (a *CommissionTierApi) DeleteCommissionTiers(c *gin.Context) {
	var req request.DeleteCommissionTiersRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := req.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := commissionTierService.DeleteCommissionTiers(req.IDs); err != nil {
		global.GVA_LOG.Error("删除阶梯等级失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("删除成功", c)
}

// UpdateCommissionTierStatus 更新等级状态
// @Tags CommissionTier
// @Summary 更新等级状态
// @Description 启用或禁用阶梯等级
// @Accept json
// @Produce json
// @Param data body request.UpdateCommissionTierStatusRequest true "状态信息"
// @Success 200 {object} response.Response "成功"
// @Router /commissionTier/updateStatus [post]
func (a *CommissionTierApi) UpdateCommissionTierStatus(c *gin.Context) {
	var req request.UpdateCommissionTierStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := req.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := commissionTierService.UpdateCommissionTierStatus(req.ID, req.Status); err != nil {
		global.GVA_LOG.Error("更新等级状态失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithMessage("更新成功", c)
}

// UpdateCommissionTierSort 批量更新等级排序
// @Tags CommissionTier
// @Summary 批量更新等级排序
// @Description 批量更新阶梯等级的排序值
// @Accept json
// @Produce json
// @Param data body request.UpdateCommissionTierSortRequest true "排序数据"
// @Success 200 {object} response.Response "成功"
// @Router /commissionTier/updateSort [post]
func (a *CommissionTierApi) UpdateCommissionTierSort(c *gin.Context) {
	var req request.UpdateCommissionTierSortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := req.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := commissionTierService.UpdateCommissionTierSort(req.Sorts); err != nil {
		global.GVA_LOG.Error("等级排序失败", zap.Error(err))
		response.FailWithMessage("等级排序失败"+err.Error(), c)
		return
	}
	response.OkWithMessage("等级排序成功", c)
}
