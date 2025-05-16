package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LoggingMiddleware 實現日誌記錄的 middleware
type LoggingMiddleware struct {
	logger Logger
}

// NewLoggingMiddleware 建立日誌 middleware
func NewLoggingMiddleware(logger Logger) Middleware {
	return &LoggingMiddleware{
		logger: logger,
	}
}

// Handle 實現 Middleware 介面
func (m *LoggingMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 產生請求 ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)

		// 記錄請求開始
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		m.logger.Info("請求開始", map[string]interface{}{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       path,
			"query":      raw,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})

		// 處理請求
		c.Next()

		// 計算處理時間
		duration := time.Since(start)

		// 記錄請求結果
		fields := map[string]interface{}{
			"request_id":  requestID,
			"status_code": c.Writer.Status(),
			"duration_ms": duration.Milliseconds(),
			"method":      c.Request.Method,
			"path":        path,
		}

		// 檢查是否有錯誤
		if len(c.Errors) > 0 {
			// 記錄錯誤
			lastError := c.Errors.Last()
			m.logger.Error("請求處理失敗", lastError, fields)
		} else {
			// 記錄成功
			m.logger.Info("請求處理完成", fields)
		}

		// 如果回應碼是錯誤，也記錄
		statusCode := c.Writer.Status()
		if statusCode >= 400 {
			fields["response_error"] = true
			m.logger.Error("請求回應錯誤", nil, fields)
		}
	}
}

// GetRequestID 從 Gin Context 獲取請求 ID (輔助函數)
func GetRequestID(c *gin.Context) string {
	requestID, exists := c.Get("request_id")
	if !exists {
		return ""
	}

	id, ok := requestID.(string)
	if !ok {
		return ""
	}

	return id
}

// LogError 記錄錯誤到 Gin Context (輔助函數)
func LogError(c *gin.Context, err error) {
	c.Error(err)
}
