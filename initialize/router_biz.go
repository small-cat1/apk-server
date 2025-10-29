package initialize

import (
	"ApkAdmin/global"
	"ApkAdmin/middleware"
	"ApkAdmin/router"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 占位方法，保证文件可以正确加载，避免go空变量检测报错，请勿删除。
func holder(routers ...*gin.RouterGroup) {
	_ = routers
	_ = router.RouterGroupApp
}

func initBizRouter(routers ...*gin.RouterGroup) {
	privateGroup := routers[0]
	publicGroup := routers[1]
	holder(publicGroup, privateGroup)
}

func initWebRouter(engine *gin.Engine) {
	PublicGroup := engine.Group("web")
	PrivateGroup := engine.Group("web")
	webRouter := router.RouterGroupApp.Web
	PublicGroup.StaticFS(global.GVA_CONFIG.Local.StorePath, justFilesFilesystem{
		http.Dir(global.GVA_CONFIG.Local.StorePath),
	})
	PrivateGroup.Use(middleware.ClientJWTAuth())
	{
		webRouter.InitBaseRouter(PublicGroup)
		webRouter.InitCategoryRouter(PublicGroup)
		webRouter.InitAppRouter(PublicGroup, PrivateGroup)
		webRouter.InitAnnouncementRouter(PublicGroup)
	}
	{
		webRouter.InitUserRouter(PrivateGroup)
		webRouter.InitOrderRouter(PrivateGroup)
		webRouter.InitMembershipPlansRoute(PrivateGroup)
		webRouter.InitPaymentRouter(PrivateGroup)
		webRouter.InitCommissionTier(PrivateGroup)
		webRouter.InitWithdrawRouter(PrivateGroup)
	}

}
