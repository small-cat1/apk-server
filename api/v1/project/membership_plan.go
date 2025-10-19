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

type MembershipPlanApi struct {
}

func (a *MembershipPlanApi) FirstMembershipPlan(c *gin.Context) {
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
	data, err := MembershipPlanService.GetMembershipPlan(id)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(data, "获取成功", c)
}

func (a *MembershipPlanApi) GetMembershipPlanList(c *gin.Context) {
	var pageInfo request2.MembershipPlanListRequest
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
	countryList, total, err := MembershipPlanService.GetMembershipPlanList(pageInfo, pageInfo.OrderKey, pageInfo.Desc)
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

func (a *MembershipPlanApi) AddMembershipPlan(c *gin.Context) {
	var req request2.MembershipPlanCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		global.GVA_LOG.Error("参数错误!", zap.Error(err))
		response.FailWithMessage("参数错误", c)
		return
	}
	// 验证创建请求
	if err := req.Validate(); err != nil {
		response.FailWithMessage(err.Error(), c)
	}
	if err := MembershipPlanService.CreateMembershipPlan(req); err != nil {
		global.GVA_LOG.Error("创建/更新失败!", zap.Error(err))
		response.FailWithMessage("创建/更新失败："+err.Error(), c)
		return
	}
	response.OkWithMessage("创建/更新成功", c)
}

func (a *MembershipPlanApi) UpdateMembershipPlan(c *gin.Context) {
	var req request2.MembershipPlanUpdateRequest
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
	err = MembershipPlanService.UpdateMembershipPlan(&req)
	if err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
		return
	}
	response.OkWithMessage("更新成功", c)
}
func (a *MembershipPlanApi) DeleteMembershipPlan(c *gin.Context) {
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
	err = MembershipPlanService.DeleteMembershipPlan(info.ID)
	if err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
		return
	}
	response.OkWithMessage("删除成功", c)
}
