// Package router 提供路由綁定的抽象介面，用來隔離具體實作（如 gin.RouterGroup）
package router

// HandlerFunc 處理請求的函式介面
type HandlerFunc interface {
	Handle(ctx any)
}

// RouterGroup 路由綁定的介面
type RouterGroup interface {
	Handle(method, path string, handler HandlerFunc)
	Group(path string) RouterGroup
}
