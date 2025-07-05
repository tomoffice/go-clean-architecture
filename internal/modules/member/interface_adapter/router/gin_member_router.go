package router

import (
	"github.com/gin-gonic/gin"
	"module-clean/internal/framework/http/gin/middleware"
	"module-clean/internal/modules/member/interface_adapter/controller"
)

type MemberRouter struct {
	controller  *controller.MemberController
	routerGroup *gin.RouterGroup
	middlewares *middleware.Container
}

func NewMemberRouter(ctrl *controller.MemberController, routerGroup *gin.RouterGroup) *MemberRouter {
	moduleRouter := routerGroup.Group("/members")
	return &MemberRouter{
		routerGroup: moduleRouter,
		controller:  ctrl,
	}
}
func (r *MemberRouter) Register() error {
	r.routerGroup.POST("", r.controller.Register)
	r.routerGroup.GET("/:id", r.controller.GetByID)
	r.routerGroup.GET("/email/:email", r.controller.GetByEmail)
	r.routerGroup.GET("", r.controller.List)
	r.routerGroup.PATCH("/:id", r.controller.UpdateProfile)
	r.routerGroup.PATCH("/:id/email", r.controller.UpdateEmail)
	r.routerGroup.PATCH("/:id/password", r.controller.UpdatePassword)
	r.routerGroup.DELETE("/:id", r.controller.Delete)
	return nil
}
