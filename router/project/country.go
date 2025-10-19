package project

import (
	"ApkAdmin/middleware"
	"github.com/gin-gonic/gin"
)

type CountryRouter struct {
}

func (r *CountryRouter) InitCountryRouter(Router *gin.RouterGroup) {
	router := Router.Group("country").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("country")
	{
		router.POST("country", countryApi.AddCountry)      // 添加国家代码
		router.PUT("country", countryApi.UpdateCountry)    // 编辑国家代码
		router.DELETE("country", countryApi.DeleteCountry) // 删除国家代码
	}
	{
		routerWithoutRecord.GET("countryList", countryApi.GetCountryList) // 国家代码列表
		routerWithoutRecord.GET("country", countryApi.FirstCountry)       // 获取单一国家代码信息

	}
}
