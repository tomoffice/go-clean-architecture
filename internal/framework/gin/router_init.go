package gin

import (
	"github.com/gin-gonic/gin"
	memberCtrl "module-clean/internal/member/interface_adapters/controller"
	memberRouter "module-clean/internal/member/interface_adapters/router"
)

func InitGinRouter(controller *memberCtrl.MemberController) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())

	adapter := NewGinRouterAdapter(engine)
	apiGroup := adapter.Group("/api/v1")
	memberRouter.RegisterMemberRoutes(apiGroup, controller)

	return engine
}
