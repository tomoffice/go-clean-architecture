package router

import (
	"github.com/gin-gonic/gin"
	"module-clean/internal/modules/member/interface_adapter/controller"
)

type MemberRouter struct {
	controller *controller.MemberController
}

func NewMemberRouter(controller *controller.MemberController) *MemberRouter {
	return &MemberRouter{controller: controller}
}

func (r *MemberRouter) RegisterRoutes(router *gin.RouterGroup) {
	member := router.Group("/members")
	member.POST("", r.controller.Register)
	member.GET("/:id", r.controller.GetByID)
	member.GET("/email/:email", r.controller.GetByEmail)
	member.GET("", r.controller.List)
	member.PATCH("/:id", r.controller.UpdateProfile)
	member.PATCH("/:id/email", r.controller.UpdateEmail)
	member.PATCH("/:id/password", r.controller.UpdatePassword)
	member.DELETE("/:id", r.controller.Delete)
}
