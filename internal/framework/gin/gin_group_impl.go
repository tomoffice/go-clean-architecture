package gin

import (
	"github.com/gin-gonic/gin"
	"module-clean/internal/shared/router"
)

type GinRootRouteGroup struct {
	engine *gin.Engine
}

func (r *GinRootRouteGroup) Group(path string) router.SubRouteGroup {
	return &GinSubRouteGroup{group: r.engine.Group(path)}
}

func NewGinRootRouteGroup(engine *gin.Engine) router.RootRouteGroup {
	return &GinRootRouteGroup{engine: engine}
}

type GinSubRouteGroup struct {
	group *gin.RouterGroup
}

func (g *GinSubRouteGroup) Handle(method, path string, handler any) {
	hf, ok := handler.(gin.HandlerFunc)
	if !ok {
		panic("handler must be gin.HandlerFunc")
	}
	g.group.Handle(method, path, hf)
}

func (g *GinSubRouteGroup) Group(path string) router.SubRouteGroup {
	return &GinSubRouteGroup{group: g.group.Group(path)}
}
