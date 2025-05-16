package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Middleware 定義 Gin 的 middleware 介面
// 這是框架層的概念，只適用於 Gin 框架
type Middleware interface {
	Handle() gin.HandlerFunc
}

// Logger 定義日誌服務介面
type Logger interface {
	Debug(msg string, fields map[string]interface{})
	Info(msg string, fields map[string]interface{})
	Warn(msg string, fields map[string]interface{})
	Error(msg string, err error, fields map[string]interface{})
}

// AuthService 定義認證服務介面
type AuthService interface {
	ValidateToken(token string) (userID string, err error)
	ValidateAPIKey(apiKey string) (userID string, err error)
	ValidateBasicAuth(username, password string) (userID string, err error)
}

// RateLimitStorage 定義速率限制儲存介面
type RateLimitStorage interface {
	GetAndIncrement(key string, window time.Duration) (count int, resetAt time.Time, err error)
	Reset(key string) error
}
