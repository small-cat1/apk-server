package project

import (
	"ApkAdmin/middleware"
	"github.com/gin-gonic/gin"
)

type CategoryRouter struct {
}

func (r *CategoryRouter) InitCategoryRouter(Router *gin.RouterGroup) {
	router := Router.Group("category").Use(middleware.OperationRecord())
	routerWithoutRecord := Router.Group("category")

	{
		router.POST("category", categoryApi.AddCategory)      // 添加分类
		router.PUT("category", categoryApi.UpdateCategory)    // 编辑分类
		router.DELETE("category", categoryApi.DeleteCategory) // 删除分类
	}
	{
		routerWithoutRecord.GET("selectCategory", categoryApi.GetSelectCategory) // 获取下拉列表分类
		routerWithoutRecord.GET("categoryList", categoryApi.GetCategoryList)     // 分类列表
		routerWithoutRecord.GET("category", categoryApi.FirstCategory)           // 获取单一分类信息
	}
}
