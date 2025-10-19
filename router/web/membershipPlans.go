package web

import "github.com/gin-gonic/gin"

type MembershipPlansRoute struct {
}

func (r *MembershipPlansRoute) InitMembershipPlansRoute(Router *gin.RouterGroup) {
	Router.GET("/membershipPlans", membershipPlansApi.GetMembershipPlans) // 获取会员套餐
}
