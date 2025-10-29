package web

import api "ApkAdmin/api/v1"

type RouterGroup struct {
	BaseRouter
	CategoryRouter
	AppRoute
	UserRouter
	OrderRouter
	MembershipPlansRoute
	PaymentRouter
	AnnouncementRouter
	CommissionTierRouter
	WithdrawRouter
	CommissionDetailRouter
}

var (
	categoryApi           = api.ApiGroupApp.WebApiGroup.CategoryApi
	appApi                = api.ApiGroupApp.WebApiGroup.AppApi
	baseApi               = api.ApiGroupApp.WebApiGroup.BaseApi
	jwtApi                = api.ApiGroupApp.WebApiGroup.JwtApi
	userApi               = api.ApiGroupApp.WebApiGroup.UserApi
	orderApi              = api.ApiGroupApp.WebApiGroup.OrderApi
	membershipPlansApi    = api.ApiGroupApp.WebApiGroup.MembershipPlansApi
	paymentApi            = api.ApiGroupApp.WebApiGroup.PaymentApi
	systemAnnouncementApi = api.ApiGroupApp.WebApiGroup.SystemAnnouncementApi
	commissionTierApi     = api.ApiGroupApp.WebApiGroup.CommissionTierApi
	withdrawApi           = api.ApiGroupApp.WebApiGroup.WithdrawApi
	commissionDetailApi   = api.ApiGroupApp.WebApiGroup.CommissionDetailApi
)
