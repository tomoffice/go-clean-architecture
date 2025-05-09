package router

import (
	"github.com/gin-gonic/gin"
)

func NewGinEngine(rootPath string, groupRegisterFuncs ...func(group *gin.RouterGroup)) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())

	ginRoot := engine.Group(rootPath)

	for _, registerFunc := range groupRegisterFuncs {
		registerFunc(ginRoot) // 呼叫註冊函式，把 ginRoot 傳進去
	}

	return engine
}