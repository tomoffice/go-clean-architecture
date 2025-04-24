package gin

import (
	"github.com/gin-gonic/gin"
	"module-clean/internal/shared/router"
)

type GinRouteGroup struct {
	group *gin.RouterGroup
}

func (g *GinRouteGroup) Handle(method, path string, handler any) {
	hf, ok := handler.(gin.HandlerFunc)
	if !ok {
		panic("handler must be gin.HandlerFunc")
	}
	g.group.Handle(method, path, hf)
}

func (g *GinRouteGroup) Group(path string) router.RouteGroup {
	return &GinRouteGroup{group: g.group.Group(path)}
}

type GinRouterAdapter struct {
	engine *gin.Engine
}

func NewGinRouterAdapter(engine *gin.Engine) router.RouterAdapter {
	return &GinRouterAdapter{engine}
}

func (a *GinRouterAdapter) Group(path string) router.RouteGroup {
	return &GinRouteGroup{group: a.engine.Group(path)}
}
