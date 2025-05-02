package router

import (
	"github.com/gin-gonic/gin"
	"module-clean/internal/member/interface_adapters/controller"
)

func NewMemberRouter(controller *controller.MemberController) func(*gin.RouterGroup) {
	return func(api *gin.RouterGroup) {
		member := api.Group("/members")

		member.POST("", controller.Register)
		member.GET("/:id", controller.GetByID)
		member.GET("/email/:email", controller.GetByEmail)
		member.GET("", controller.List)
		member.PUT("/:id", controller.Update)
		member.DELETE("/:id", controller.Delete)
	}
}
