package member

import (
	"github.com/gin-gonic/gin"
	"module-clean/internal/modules/member/interface_adapter/router"
)

type Module struct {
	moduleRouter *router.MemberRouter
}

func NewModule(moduleRouter *router.MemberRouter) *Module {
	return &Module{

		moduleRouter: moduleRouter,
	}
}
func (m *Module) RegisterRoutes(routerGroup *gin.RouterGroup) {
	m.moduleRouter.RegisterRoutes(routerGroup)
}
