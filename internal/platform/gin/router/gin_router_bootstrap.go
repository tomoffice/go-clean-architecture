package router

import (
	"github.com/gin-gonic/gin"
	sharedrouter "module-clean/internal/shared/router"
)

func NewGinEngine(rootPath string, groupRegisterFunc ...func(group sharedrouter.RouterGroup)) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())

	ginRoot := NewGinRouterGroup(engine.Group(rootPath))

	for _, registerFunc := range groupRegisterFunc {
		registerFunc(ginRoot)
	}

	return engine
}
