package router

import (
	"github.com/gin-gonic/gin"
)

// GinHandler 將 gin.HandlerFunc 包裝成 shared.HandlerFunc
type GinHandler struct {
	H gin.HandlerFunc
}

func NewGinHandler(h gin.HandlerFunc) *GinHandler 
func NewGinHandler(h gin.HandlerFunc) *GinHandler {
	return &GinHandler{H: h}
}
func (g *GinHandler) Handle(ctx any) {
	// 保證 ctx 是 *gin.Context
	g.H(ctx.(*gin.Context))
}
