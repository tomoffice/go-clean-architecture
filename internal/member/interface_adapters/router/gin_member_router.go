package router

import (
	"github.com/gin-gonic/gin"
	"module-clean/internal/member/interface_adapters/controller"
	"module-clean/internal/shared/router"
)

func RegisterMemberRoutes(rg router.RouteGroup, controller *controller.MemberController) {
	member := rg.Group("/members")
	member.Handle("POST", "", gin.HandlerFunc(controller.Register))
	member.Handle("GET", "", gin.HandlerFunc(controller.List))
	member.Handle("GET", "/:id", gin.HandlerFunc(controller.GetByID))
	member.Handle("PUT", "/:id", gin.HandlerFunc(controller.Update))
	member.Handle("DELETE", "/:id", gin.HandlerFunc(controller.Delete))
}
