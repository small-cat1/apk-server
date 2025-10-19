package web

import (
	"ApkAdmin/global"
	"ApkAdmin/model/common/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type MembershipPlansApi struct {
}

// GetMembershipPlans 获取会员套餐
func (a MembershipPlansApi) GetMembershipPlans(c *gin.Context) {
	plans, err := membershipPlanService.GetAllMembershipPlan()
	if err != nil {
		global.GVA_LOG.Error("获取会员套餐列表失败", zap.Error(err))
		response.OkWithMessage("获取会员套餐列表失败", c)
		return
	}
	response.OkWithDetailed(plans, "success", c)
	return
}
