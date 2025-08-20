package router

import (
	"github.com/gin-gonic/gin"
	ginadapter "github.com/tomoffice/go-clean-architecture/internal/framework/http/gin/adapter"
	memberhttp "github.com/tomoffice/go-clean-architecture/internal/interface_adapter/transport/http"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/interface_adapter/controller"
)

type MemberRouter struct {
	controller *controller.MemberController
	router     memberhttp.Router
}

func NewMemberRouter(ctrl *controller.MemberController, routerGroup *gin.RouterGroup) *MemberRouter {
	moduleGroup := routerGroup.Group("/members")
	return &MemberRouter{
		controller: ctrl,
		router:     ginadapter.NewRouter(moduleGroup), // ← 這裡建立抽象 router
	}
}

func (r *MemberRouter) Register() error {
	r.router.POST("", r.controller.Register)
	r.router.GET("/:id", r.controller.GetByID)
	r.router.GET("/email/:email", r.controller.GetByEmail)
	r.router.GET("", r.controller.List)
	r.router.PATCH("/:id", r.controller.UpdateProfile)
	r.router.PATCH("/:id/email", r.controller.UpdateEmail)
	r.router.PATCH("/:id/password", r.controller.UpdatePassword)
	r.router.DELETE("/:id", r.controller.Delete)
	return nil
}
