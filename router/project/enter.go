package project

import api "ApkAdmin/api/v1"

type RouterGroup struct {
	CountryRouter
	CategoryRouter
	AppPackageRouter
	MembershipPlanRouter
	UserRouter
	ApplicationRouter
	MembershipOrderRouter
	PaymentProviderRouter
	PaymentAccountRouter
	WebsiteConfigRouter
	AppAccountRouter
	SystemAnnouncementRouter
	CommissionTierRouter
}

var (
	categoryApi           = api.ApiGroupApp.ProjectApiGroup.CategoryApi
	countryApi            = api.ApiGroupApp.ProjectApiGroup.CountryApi
	appPackageApi         = api.ApiGroupApp.ProjectApiGroup.AppPackageApi
	membershipPlanApi     = api.ApiGroupApp.ProjectApiGroup.MembershipPlanApi
	userApi               = api.ApiGroupApp.ProjectApiGroup.UserApi
	applicationApi        = api.ApiGroupApp.ProjectApiGroup.ApplicationApi
	membershipOrderApi    = api.ApiGroupApp.ProjectApiGroup.MembershipOrderApi
	paymentProviderApi    = api.ApiGroupApp.ProjectApiGroup.PaymentProviderApi
	paymentAccountApi     = api.ApiGroupApp.ProjectApiGroup.PaymentAccountApi
	websiteConfigApi      = api.ApiGroupApp.ProjectApiGroup.WebsiteConfigApi
	appAccountApi         = api.ApiGroupApp.ProjectApiGroup.AppAccountApi
	systemAnnouncementApi = api.ApiGroupApp.ProjectApiGroup.SystemAnnouncementApi
	commissionTierApi     = api.ApiGroupApp.ProjectApiGroup.CommissionTierApi
)
