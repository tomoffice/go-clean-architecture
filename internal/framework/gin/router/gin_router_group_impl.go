package router

import (
	"github.com/gin-gonic/gin"
	sharedrouter "module-clean/internal/shared/router"
)

type GinRouterGroup struct {
	group *gin.RouterGroup
}

func NewGinRouterGroup(group *gin.RouterGroup) sharedrouter.RouterGroup {
	return &GinRouterGroup{group: group}
}
func (g *GinRouterGroup) Handle(method, path string, handler sharedrouter.HandlerFunc) {
	g.group.Handle(method, path, func(c *gin.Context) {
		handler.Handle(c)
	})
}

func (g *GinRouterGroup) Group(path string) sharedrouter.RouterGroup {
	return &GinRouterGroup{group: g.group.Group(path)}
}
