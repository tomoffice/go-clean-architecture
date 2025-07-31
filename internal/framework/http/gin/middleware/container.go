package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/tomoffice/go-clean-architecture/internal/framework/http/gin/middleware/cors"
	"github.com/tomoffice/go-clean-architecture/internal/framework/http/gin/middleware/logging"
	"github.com/tomoffice/go-clean-architecture/pkg/logger"
	"github.com/tomoffice/go-clean-architecture/pkg/tracer"
)

// Container 提供所有可用的middleware
type Container struct {
	logger    logger.Logger
	tracer    tracer.Tracer
	cors      gin.HandlerFunc
	auth      gin.HandlerFunc
	rateLimit gin.HandlerFunc
	logging   gin.HandlerFunc
	// 可以繼續添加其他middleware
}

// NewContainer 創建並初始化middleware容器
func NewContainer(logger logger.Logger, tracer tracer.Tracer) *Container {
	return &Container{
		logger:  logger,
		tracer:  tracer,
		cors:    cors.NewCORSMiddleware(cors.DefaultCORSConfig()).HandlerFunc(),
		logging: logging.NewLoggingMiddleware(logging.DefaultLoggingConfig(), logger, tracer).HandlerFunc(),
		// 暫時將其他中間件設為nil，之後實現時再添加
		auth:      nil,
		rateLimit: nil,
	}
}

// CORS 返回CORS中間件，如果不存在則返回nil
func (c *Container) CORS() gin.HandlerFunc {
	return c.cors
}

// Auth 返回認證中間件，如果不存在則返回nil
func (c *Container) Auth() gin.HandlerFunc {
	return c.auth
}

// RateLimit 返回限流中間件，如果不存在則返回nil
func (c *Container) RateLimit() gin.HandlerFunc {
	return c.rateLimit
}

// HasCORS 檢查是否有CORS中間件
func (c *Container) HasCORS() bool { return c.cors != nil }

// HasAuth 檢查是否有認證中間件
func (c *Container) HasAuth() bool { return c.auth != nil }

// HasRateLimit 檢查是否有限流中間件
func (c *Container) HasRateLimit() bool { return c.rateLimit != nil }

// Logging 返回日誌中間件，如果不存在則返回nil
func (c *Container) Logging() gin.HandlerFunc {
	return c.logging
}

// HasLogging 檢查是否有日誌中間件
func (c *Container) HasLogging() bool { return c.logging != nil }
