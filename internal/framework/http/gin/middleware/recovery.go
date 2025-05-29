package middleware

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware 實現 panic 恢復的 middleware
type RecoveryMiddleware struct {
	logger       Logger
	stack        bool
	stackSize    int
	printStack   bool
	errorHandler RecoveryErrorHandler
}

// RecoveryConfig Recovery middleware 配置
type RecoveryConfig struct {
	Logger       Logger               // 日誌器
	Stack        bool                 // 是否記錄堆疊資訊
	StackSize    int                  // 堆疊大小（預設: 4KB）
	PrintStack   bool                 // 是否列印堆疊到控制台
	ErrorHandler RecoveryErrorHandler // 自定義錯誤處理器
}

// RecoveryErrorHandler 定義 panic 錯誤處理器
type RecoveryErrorHandler func(c *gin.Context, err interface{})

// DefaultRecoveryConfig 提供預設的 Recovery 配置
func DefaultRecoveryConfig() RecoveryConfig {
	return RecoveryConfig{
		Stack:      true,
		StackSize:  4 << 10, // 4 KB
		PrintStack: gin.IsDebugging(),
	}
}

// NewRecoveryMiddleware 建立 panic 恢復 middleware
func NewRecoveryMiddleware(config RecoveryConfig) Middleware {
	if config.StackSize == 0 {
		config.StackSize = 4 << 10 // 4 KB
	}

	if config.ErrorHandler == nil {
		config.ErrorHandler = defaultErrorHandler
	}

	return &RecoveryMiddleware{
		logger:       config.Logger,
		stack:        config.Stack,
		stackSize:    config.StackSize,
		printStack:   config.PrintStack,
		errorHandler: config.ErrorHandler,
	}
}

// Handle 實現 Middleware 介面
func (m *RecoveryMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 檢查連接是否中斷
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") ||
							strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				// 記錄錯誤
				if m.logger != nil {
					stack := ""
					if m.stack {
						stack = string(debug.Stack())
					}

					httpRequest, _ := httputil.DumpRequest(c.Request, false)
					headers := strings.Split(string(httpRequest), "\r\n")
					for idx, header := range headers {
						current := strings.Split(header, ":")
						if current[0] == "Authorization" {
							headers[idx] = current[0] + ": *"
						}
					}

					fields := map[string]interface{}{
						"time":       time.Now(),
						"error":      err,
						"dto":        strings.Join(headers, "\r\n"),
						"stack":      stack,
						"request_id": c.GetString("request_id"),
					}

					if brokenPipe {
						fields["broken_pipe"] = true
					}

					m.logger.Error("恢復 panic", fmt.Errorf("%v", err), fields)
				}

				// 列印堆疊到控制台
				if m.printStack {
					fmt.Fprintf(gin.DefaultErrorWriter, "[Recovery] %s panic recovered:\n%s\n%s%s\n",
						time.Now().Format("2006/01/02 - 15:04:05"),
						err,

						string(debug.Stack()),
					)
				}

				// 如果是連接中斷，返回錯誤但不嘗試寫入回應
				if brokenPipe {
					// 如果連接已中斷，嘗試寫入回應可能會再次 panic
					c.Error(err.(error))
					c.Abort()
				} else {
					// 使用自定義錯誤處理器
					m.errorHandler(c, err)
				}
			}
		}()

		c.Next()
	}
}

// defaultErrorHandler 預設的錯誤處理器
func defaultErrorHandler(c *gin.Context, err interface{}) {
	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"error":      "伺服器內部錯誤",
		"message":    fmt.Sprintf("%v", err),
		"status":     http.StatusInternalServerError,
		"path":       c.Request.URL.Path,
		"method":     c.Request.Method,
		"request_id": c.GetString("request_id"),
		"timestamp":  time.Now().Format(time.RFC3339),
	})
}

// RecoveryWithWriter 建立寫入到特定 writer 的 Recovery middleware
func RecoveryWithWriter(out gin.ResponseWriter, recovery ...RecoveryErrorHandler) gin.HandlerFunc {
	if len(recovery) > 0 {
		return NewRecoveryMiddleware(RecoveryConfig{
			Stack:        true,
			PrintStack:   true,
			ErrorHandler: recovery[0],
		}).Handle()
	}

	return NewRecoveryMiddleware(DefaultRecoveryConfig()).Handle()
}

// CustomRecoveryConfig 提供自定義的 Recovery 配置
func CustomRecoveryConfig(logger Logger, handler RecoveryErrorHandler) RecoveryConfig {
	return RecoveryConfig{
		Logger:       logger,
		Stack:        true,
		StackSize:    4 << 10,
		PrintStack:   gin.IsDebugging(),
		ErrorHandler: handler,
	}
}

// PanicInfo 用於記錄 panic 資訊的結構
type PanicInfo struct {
	Time      time.Time
	Error     interface{}
	Stack     string
	RequestID string
	Path      string
	Method    string
	ClientIP  string
}

// GetPanicInfo 從 context 中提取 panic 資訊
func GetPanicInfo(c *gin.Context, err interface{}) *PanicInfo {
	return &PanicInfo{
		Time:      time.Now(),
		Error:     err,
		Stack:     string(debug.Stack()),
		RequestID: c.GetString("request_id"),
		Path:      c.Request.URL.Path,
		Method:    c.Request.Method,
		ClientIP:  c.ClientIP(),
	}
}
