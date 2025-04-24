// Package router 提供路由綁定的抽象介面，用來隔離具體實作（如 gin.RouterGroup）
package router

// RouterAdapter 應用入口的 router glue 封裝（通常用於 /api/v1 起始點）
// 可對應 gin.Engine.Group(path) 或 fiber.App.Group(path)
type RouterAdapter interface {
	Group(path string) RouteGroup
}

// RouteGroup 提供模組級別的路由設定（如 /members、/users）
// 所有 handler 綁定都走這層，以確保不耦合底層框架
type RouteGroup interface {
	Handle(method string, path string, handler any)
	Group(path string) RouteGroup
}
