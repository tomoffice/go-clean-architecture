package logging

import (
	"bytes"
	"context"
	"io"
	"strings"
	"time"
	
	"github.com/gin-gonic/gin"
	"github.com/tomoffice/go-clean-architecture/pkg/logger"
	"github.com/tomoffice/go-clean-architecture/pkg/tracer"
	"go.opentelemetry.io/otel/trace"
)

// LoggingConfig 定義 logging middleware 的配置
type LoggingConfig struct {
	// SkipPaths 定義要跳過追蹤的路徑
	SkipPaths []string
	
	// LogRequestBody 是否記錄請求 body
	LogRequestBody bool
	
	// LogResponseBody 是否記錄回應 body
	LogResponseBody bool
	
	// MaxBodySize body 記錄的最大大小（bytes），超過此大小將被截斷
	MaxBodySize int
	
	// LoggableContentTypes 可記錄的 Content-Type 白名單
	// 空列表表示記錄所有類型
	LoggableContentTypes []string
}

// DefaultLoggingConfig 返回預設配置
func DefaultLoggingConfig() LoggingConfig {
	return LoggingConfig{
		SkipPaths:            []string{"/health", "/metrics"},
		LogRequestBody:       true,
		LogResponseBody:      true,
		MaxBodySize:          1024, // 1KB
		LoggableContentTypes: []string{"application/json", "text/plain", "application/xml"},
	}
}

// LoggingMiddleware 是 Gin 中間件，負責 HTTP 請求的日誌記錄
// 使用 tracer 生成 trace ID 用於日誌關聯和 HTTP headers
// 注意：這不是完整的 observability 解決方案，真正的 tracing 應使用專門的 middleware
type LoggingMiddleware struct {
	config LoggingConfig
	logger logger.Logger
	tracer tracer.Tracer
}

// NewLoggingMiddleware 建立新的 logging middleware
func NewLoggingMiddleware(config LoggingConfig, logger logger.Logger, tracer tracer.Tracer) *LoggingMiddleware {
	return &LoggingMiddleware{
		config: config,
		logger: logger,
		tracer: tracer,
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

		// 記錄開始時間
		startTime := time.Now()

		// 使用 tracer 創建 span（主要用於生成 trace ID）
		ctx, span := lm.tracer.Start(c.Request.Context(), "http-request")
		defer span.End()

		// 創建帶有 trace context 的 logger
		contextLogger := lm.logger.WithContext(ctx)

		// 讀取 request body
		requestBody := lm.readRequestBody(c)

		// 建立日誌欄位
		requestFields := []logger.Field{
			logger.NewField("method", c.Request.Method),
			logger.NewField("path", c.Request.URL.Path),
			logger.NewField("remote_addr", c.ClientIP()),
			logger.NewField("user_agent", c.Request.UserAgent()),
		}
		
		// 如果有 request body，加入到日誌中
		if requestBody != "" {
			requestFields = append(requestFields, logger.NewField("request_body", requestBody))
		}

		// 記錄請求開始
		contextLogger.Info("HTTP 請求開始", requestFields...)

		// 將更新後的 context 設回 request
		c.Request = c.Request.WithContext(ctx)

		// 設置 trace ID 到 response header 用於請求追蹤
		lm.setTraceHeaders(c, ctx)

		// 如果需要記錄 response body，設置自定義 writer
		var responseBodyWriter *ResponseBodyWriter
		if lm.config.LogResponseBody {
			responseBodyWriter = &ResponseBodyWriter{
				ResponseWriter: c.Writer,
				body:          bytes.NewBufferString(""),
			}
			c.Writer = responseBodyWriter
		}

		// 處理請求
		c.Next()

		// 計算 latency
		latency := time.Since(startTime)
		statusCode := c.Writer.Status()
		responseSize := c.Writer.Size()

		// 建立回應日誌欄位
		responseFields := []logger.Field{
			logger.NewField("status_code", statusCode),
			logger.NewField("response_size", responseSize),
			logger.NewField("latency_ms", latency.Milliseconds()),
		}

		// 如果有 response body，加入到日誌中
		if lm.config.LogResponseBody && responseBodyWriter != nil {
			responseBody := responseBodyWriter.body.String()
			if len(responseBody) > lm.config.MaxBodySize {
				responseBody = responseBody[:lm.config.MaxBodySize] + "...[truncated]"
			}
			responseFields = append(responseFields, logger.NewField("response_body", responseBody))
		}

		// 記錄請求完成
		contextLogger.Info("HTTP 請求完成", responseFields...)
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

// shouldLogBody 檢查是否應該記錄 body
func (lm *LoggingMiddleware) shouldLogBody(contentType string) bool {
	if len(lm.config.LoggableContentTypes) == 0 {
		return true // 空白名單表示記錄所有類型
	}
	
	for _, loggableType := range lm.config.LoggableContentTypes {
		if strings.Contains(strings.ToLower(contentType), strings.ToLower(loggableType)) {
			return true
		}
	}
	return false
}

// readRequestBody 讀取並記錄 request body
func (lm *LoggingMiddleware) readRequestBody(c *gin.Context) string {
	if !lm.config.LogRequestBody {
		return ""
	}
	
	contentType := c.GetHeader("Content-Type")
	if !lm.shouldLogBody(contentType) {
		return "[unsupported content type]"
	}
	
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		// 重新設置 body 供後續使用
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}
	
	// 檢查大小限制
	if len(bodyBytes) > lm.config.MaxBodySize {
		return string(bodyBytes[:lm.config.MaxBodySize]) + "...[truncated]"
	}
	
	return string(bodyBytes)
}

// ResponseBodyWriter 自定義 ResponseWriter 用於攔截 response body
type ResponseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 實作 ResponseWriter.Write 方法
func (w *ResponseBodyWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// setTraceHeaders 設置 trace 資訊到 response header
func (lm *LoggingMiddleware) setTraceHeaders(c *gin.Context, ctx context.Context) {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return
	}
	
	spanCtx := span.SpanContext()
	if spanCtx.IsValid() {
		c.Header("X-Trace-Id", spanCtx.TraceID().String())
		c.Header("X-Span-Id", spanCtx.SpanID().String())
	}
}