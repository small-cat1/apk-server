package project

import (
	"ApkAdmin/middleware"
	"github.com/gin-gonic/gin"
)

type MembershipPlanRouter struct {
}

func (r *MembershipPlanRouter) InitMembershipPlanRouter(Router *gin.RouterGroup) {
	router := Router.Group("membershipPlan").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("membershipPlan")
	{
		router.POST("createMembershipPlan", membershipPlanApi.AddMembershipPlan)      // 添加会员套餐
		router.PUT("updateMembershipPlan", membershipPlanApi.UpdateMembershipPlan)    // 编辑会员套餐
		router.DELETE("deleteMembershipPlan", membershipPlanApi.DeleteMembershipPlan) // 删除会员套餐
	}
	{
		routerWithoutRecord.GET("getMembershipPlanList", membershipPlanApi.GetMembershipPlanList) // 会员套餐列表
		routerWithoutRecord.GET("findMembershipPlan", membershipPlanApi.FirstMembershipPlan)      // 获取单一会员套餐信息
	}
}
