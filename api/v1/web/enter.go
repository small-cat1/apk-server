package web

import "ApkAdmin/service"

type ApiGroup struct {
	CategoryApi
	AppApi
	BaseApi
	JwtApi
	UserApi
	OrderApi
	MembershipPlansApi
	PaymentApi
	SystemAnnouncementApi
	CommissionTierApi
	WithdrawApi
	CommissionDetailApi
}

var (
	CategoryService           = service.ServiceGroupApp.ProjectServiceGroup.CategoryService
	AppService                = service.ServiceGroupApp.ProjectServiceGroup.ApplicationService
	UserService               = service.ServiceGroupApp.ProjectServiceGroup.UserService
	jwtService                = service.ServiceGroupApp.SystemServiceGroup.JwtService
	websiteConfigService      = service.ServiceGroupApp.ProjectServiceGroup.SystemConfigService
	userMembershipService     = service.ServiceGroupApp.ProjectServiceGroup.UserMembershipService
	membershipPlanService     = service.ServiceGroupApp.ProjectServiceGroup.MembershipPlanService
	paymentProviderService    = service.ServiceGroupApp.ProjectServiceGroup.PaymentProviderService
	systemAnnouncementService = service.ServiceGroupApp.ProjectServiceGroup.SystemAnnouncementService
	commissionTierService     = service.ServiceGroupApp.ProjectServiceGroup.CommissionTierService
	systemConfigService       = service.ServiceGroupApp.ProjectServiceGroup.SystemConfigService
	commissionDetailService   = service.ServiceGroupApp.ProjectServiceGroup.CommissionDetailService
)
