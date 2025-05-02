package gin

import (
	"github.com/gin-gonic/gin"
)

func InitGinRouter(rootPath string, routerGroups ...func(*gin.RouterGroup)) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())

	api := engine.Group(rootPath)
	for _, group := range routerGroups {
		group(api)
	}

	return engine
}
