package web

import (
	"ApkAdmin/global"
	"ApkAdmin/model/common/response"
	"ApkAdmin/model/project/request"
	"ApkAdmin/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CommissionDetailApi struct{}

// GetCommissionDetailList 分页获取分佣明细列表
// @Tags CommissionDetail
// @Summary 分页获取分佣明细列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.CommissionDetailSearch true "分页获取分佣明细列表"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "分页获取分佣明细列表,返回包括列表,总数,页码,每页数量"
// @Router /commissionDetail/getCommissionDetailList [get]
func (api *CommissionDetailApi) GetCommissionDetailList(c *gin.Context) {
	var search request.ClientCommissionDetailSearch
	err := c.ShouldBindQuery(&search)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 设置默认分页参数
	if search.Page < 1 {
		search.Page = 1
	}
	if search.PageSize < 1 {
		search.PageSize = 10
	}
	// 从上下文中获取用户ID（如果需要根据登录用户查询）
	userId := utils.GetUserID(c)
	list, err := commissionDetailService.ClientGetCommissionDetailList(userId, search)
	if err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
		return
	}
	response.OkWithDetailed(list, "获取成功", c)
}
