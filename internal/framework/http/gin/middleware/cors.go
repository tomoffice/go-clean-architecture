package middleware

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware 實現跨域資源共享的 middleware
type CORSMiddleware struct {
	config CORSConfig
}

// CORSConfig 定義 CORS 參數
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// DefaultCORSConfig 提供預設的 CORS 配置
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}

// NewCORSMiddleware 建構 CORS 中介層
func NewCORSMiddleware(config CORSConfig) Middleware {
	// 如果沒有設定，使用預設值
	if len(config.AllowOrigins) == 0 {
		config.AllowOrigins = []string{"*"}
	}
	if len(config.AllowMethods) == 0 {
		config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	}

	return &CORSMiddleware{config: config}
}

// Handle 實現 Middleware 介面
func (m *CORSMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// 檢查是否允許的來源
		if m.isOriginAllowed(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
		} else if m.isWildcardAllowed() {
			c.Header("Access-Control-Allow-Origin", "*")
		}

		// 設置其他 CORS headers
		c.Header("Access-Control-Allow-Methods", strings.Join(m.config.AllowMethods, ","))
		c.Header("Access-Control-Allow-Headers", strings.Join(m.config.AllowHeaders, ","))

		if len(m.config.ExposeHeaders) > 0 {
			c.Header("Access-Control-Expose-Headers", strings.Join(m.config.ExposeHeaders, ","))
		}

		if m.config.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		c.Header("Access-Control-Max-Age", strconv.FormatInt(int64(m.config.MaxAge.Seconds()), 10))

		// 處理預檢請求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// isOriginAllowed 檢查來源是否被允許
func (m *CORSMiddleware) isOriginAllowed(origin string) bool {
	for _, allowed := range m.config.AllowOrigins {
		if allowed == origin {
			return true
		}
	}
	return false
}

// isWildcardAllowed 檢查是否允許所有來源
func (m *CORSMiddleware) isWildcardAllowed() bool {
	for _, allowed := range m.config.AllowOrigins {
		if allowed == "*" {
			return true
		}
	}
	return false
}
