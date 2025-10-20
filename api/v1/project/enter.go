package project

import "ApkAdmin/service"

type ApiGroup struct {
	CountryApi
	CategoryApi
	AppPackageApi
	MembershipPlanApi
	UserApi
	ApplicationApi
	MembershipOrderApi
	PaymentProviderApi
	PaymentAccountApi
	WebsiteConfigApi
	AppAccountApi
	SystemAnnouncementApi
	CommissionTierApi
}

var (
	CategoryService              = service.ServiceGroupApp.ProjectServiceGroup.CategoryService
	CountryService               = service.ServiceGroupApp.ProjectServiceGroup.CountryService
	AppPackageService            = service.ServiceGroupApp.ProjectServiceGroup.AppPackageService
	AppAccountService            = service.ServiceGroupApp.ProjectServiceGroup.AppAccountService
	MembershipPlanService        = service.ServiceGroupApp.ProjectServiceGroup.MembershipPlanService
	UserService                  = service.ServiceGroupApp.ProjectServiceGroup.UserService
	ApplicationService           = service.ServiceGroupApp.ProjectServiceGroup.ApplicationService
	membershipOrderService       = service.ServiceGroupApp.ProjectServiceGroup.MembershipOrderService
	membershipOrderRefundService = service.ServiceGroupApp.ProjectServiceGroup.MembershipOrderRefundService
	PaymentProviderService       = service.ServiceGroupApp.ProjectServiceGroup.PaymentProviderService
	paymentAccountService        = service.ServiceGroupApp.ProjectServiceGroup.PaymentAccountService
	websiteConfigService         = service.ServiceGroupApp.ProjectServiceGroup.SystemConfigService
	systemAnnouncementService    = service.ServiceGroupApp.ProjectServiceGroup.SystemAnnouncementService
	sysUserService               = service.ServiceGroupApp.SystemServiceGroup.UserService
	commissionTierService        = service.ServiceGroupApp.ProjectServiceGroup.CommissionTierService
)
