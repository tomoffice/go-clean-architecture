// Package router 提供路由綁定的抽象介面，用來隔離具體實作（如 gin.RouterGroup）
package router

// RootRouteGroup 應用入口的 router glue 封裝（通常用於 /api/v1 起始點）
// 可對應 gin.Engine.Group(path) 或 fiber.App.Group(path)
type RootRouteGroup interface {
	Group(path string) SubRouteGroup
}

// SubRouteGroup 提供模組級別的路由設定（如 /members、/users）
// 所有 handler 綁定都走這層，以確保可以抽換框架
type SubRouteGroup interface {
	Handle(method string, path string, handler any)
	Group(path string) SubRouteGroup
}
