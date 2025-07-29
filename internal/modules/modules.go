package modules

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

// Module 模組接口 - 產品
type Module interface {
	// Name 獲取模組名稱
	Name() string
	// Setup 初始化模組（註冊路由等）
	Setup() error
}

// ModuleFactory 模組工廠接口 - 工廠方法模式的核心
type ModuleFactory interface {
	// CreateModule 創建並返回模組實例
	CreateModule(db *sqlx.DB, rg *gin.RouterGroup) (Module, error)
}
