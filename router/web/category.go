package web

import (
	"github.com/gin-gonic/gin"
)

type CategoryRouter struct {
}

func (r *CategoryRouter) InitCategoryRouter(Router *gin.RouterGroup) {
	Router.GET("getTrendingCategory", categoryApi.GetTrendingCategory)   // 获取热门分类
	Router.GET("categories", categoryApi.FindCategory)                   //获取所有分类
	Router.GET("account/categories", categoryApi.FindAccountAppCategory) //获取账号列表的分类
}
