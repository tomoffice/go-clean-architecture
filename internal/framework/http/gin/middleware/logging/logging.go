package logging

import (
	"github.com/gin-gonic/gin"
	"github.com/tomoffice/go-clean-architecture/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// LoggingConfig 定義 logging middleware 的配置
type LoggingConfig struct {
	// SkipPaths 定義要跳過追蹤的路徑
	SkipPaths []string
	// TracerName OpenTelemetry tracer 名稱
	TracerName string
}

// DefaultLoggingConfig 返回預設配置
func DefaultLoggingConfig() LoggingConfig {
	return LoggingConfig{
		SkipPaths:  []string{"/health", "/metrics"},
		TracerName: "gin-http-server",
	}
}

// LoggingMiddleware 是 Gin 中間件，使用 OpenTelemetry 自動創建 span
// 並將其存入 context 中，讓後續的 logger 可以使用
type LoggingMiddleware struct {
	config LoggingConfig
	logger logger.Logger
}

// NewLoggingMiddleware 建立新的 logging middleware
func NewLoggingMiddleware(config LoggingConfig, logger logger.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{
		config: config,
		logger: logger,
	}
}

// HandlerFunc 返回 Gin 中間件函數
func (lm *LoggingMiddleware) HandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 檢查是否要跳過此路徑
		if lm.shouldSkip(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 按需創建 HTTP tracer
		httpTracer := otel.Tracer("http-request")
		ctx := c.Request.Context()
		ctx, span := httpTracer.Start(ctx, c.Request.Method+" "+c.Request.URL.Path)

		// 使用注入的 logger 記錄請求開始
		requestLogger := lm.logger.WithContext(ctx)
		requestLogger.Info("HTTP 請求開始",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("remote_addr", c.ClientIP()),
		)

		defer func() {
			// 記錄 HTTP 回應資訊
			span.SetAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.path", c.Request.URL.Path),
				attribute.Int("http.status_code", c.Writer.Status()),
				attribute.Int("http.response_size", c.Writer.Size()),
			)

			// 使用注入的 logger 記錄請求完成
			requestLogger.Info("HTTP 請求完成",
				zap.Int("status", c.Writer.Status()),
				zap.Int("response_size", c.Writer.Size()),
			)

			span.End()
		}()

		// 將更新後的 context 設回 request
		c.Request = c.Request.WithContext(ctx)

		// 將 trace 資訊加入 response header，方便 client 追蹤
		spanCtx := span.SpanContext()
		if spanCtx.IsValid() {
			c.Header("X-Trace-Id", spanCtx.TraceID().String())
			c.Header("X-Span-Id", spanCtx.SpanID().String())
		}

		c.Next()
	}
}

// shouldSkip 檢查是否要跳過此路徑
func (lm *LoggingMiddleware) shouldSkip(path string) bool {
	for _, skipPath := range lm.config.SkipPaths {
		if path == skipPath {
			return true
		}
	}
	return false
}
