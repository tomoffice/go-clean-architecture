package router

import (
	router2 "module-clean/internal/framework/http/gin/router"
	"module-clean/internal/modules/member/interface_adapter/controller"
	sharedrouter "module-clean/internal/shared/interface_adapter/router"
)

type MemberRouter struct {
	controller *controller.MemberController
}

func NewMemberRouter(controller *controller.MemberController) *MemberRouter {
	return &MemberRouter{controller: controller}
}

func (r *MemberRouter) RegisterRoutes(router sharedrouter.RouterGroup) {
	member := router.Group("/members")
	member.Handle("POST", "", router2.NewGinHandler(r.controller.Register))
	member.Handle("GET", "/:id", router2.NewGinHandler(r.controller.GetByID))
	member.Handle("GET", "/email/:email", router2.NewGinHandler(r.controller.GetByEmail))
	member.Handle("GET", "", router2.NewGinHandler(r.controller.List))
	member.Handle("PUT", "/:id", router2.NewGinHandler(r.controller.Update))
	member.Handle("DELETE", "/:id", router2.NewGinHandler(r.controller.Delete))
}
