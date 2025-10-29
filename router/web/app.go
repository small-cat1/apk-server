package web

import "github.com/gin-gonic/gin"

type AppRoute struct {
}

func (r *AppRoute) InitAppRouter(PublicRouter *gin.RouterGroup, PrivateRoute *gin.RouterGroup) {
	{
		PublicRouter.GET("getHotOrRecommendApplications", appApi.ListHotOrRecommendApp) // 获取热门和推荐的应用
		PublicRouter.GET("categories/apps", appApi.GetFilterApps)                       //条件查找分类列表下的应用
		PublicRouter.GET("accounts/apps", appApi.GetAccountAppsListByCategory)          //根据分类获取应用账号
		PublicRouter.GET("app/searchApp", appApi.SearchApps)                            //搜索应用
	}
	{
		PrivateRoute.POST("app/downloadApp", appApi.DownloadApp) //下载应用

	}
}
