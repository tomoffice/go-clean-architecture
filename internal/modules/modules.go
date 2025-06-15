package modules

import "github.com/gin-gonic/gin"

type Module interface {
	RegisterRoutes(routerGroup *gin.RouterGroup)
}
