package modules

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/tomoffice/go-clean-architecture/pkg/logger"
	"github.com/tomoffice/go-clean-architecture/pkg/tracer"
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
	// CreateModule 創建並返回模組實例，注入 logger 和 tracer 以支援業務層日誌記錄和追蹤
	CreateModule(db *sqlx.DB, rg *gin.RouterGroup, logger logger.Logger, tracer tracer.Tracer) (Module, error)
}
