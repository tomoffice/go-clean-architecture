package cors

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CORSConfig 配置跨域資源共享(CORS)中間件的選項
//
// 所有字段說明：
//   - AllowOrigins: 允許的來源域名列表，例如 ["https://example.com", "https://api.example.com"]
//   - AllowMethods: 允許的 HTTP 方法列表，例如 ["GET", "POST", "PUT"]
//   - AllowHeaders: 允許的 HTTP 請求頭列表，例如 ["Content-Type", "Authorization"]
//   - ExposeHeaders: 允許瀏覽器訪問的響應頭列表，例如 ["Content-Length", "X-Request-ID"]
//   - AllowCredentials: 是否允許請求包含認證信息（cookies、HTTP認證、客戶端SSL證書）
//   - MaxAge: 預檢請求結果的快取時間（秒），即 OPTIONS 請求的有效期
//
// 安全注意事項：
//  1. 如果 AllowOrigins 包含通配符 "*"，則 AllowCredentials 必須設為 false
//  2. 生產環境建議明確列出允許的來源域名，避免使用通配符 "*"
type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
	MaxAge           int
}

// NewCORSConfig 創建自訂的 CORS 配置
//
// 參數說明:
//   - allowOrigins: 允許的來源域名列表，例如：["https://example.com", "https://api.com"]
//   - allowMethods: 允許的 HTTP 方法列表，例如：["GET", "POST", "PUT"]
//   - allowHeaders: 允許的請求頭列表，例如：["Content-Type", "Authorization"]
//   - exposeHeaders: 允許瀏覽器訪問的響應頭列表，例如：["Content-Length"]
//   - allowCredentials: 是否允許包含認證信息的請求（cookies、HTTP認證等）
//   - maxAge: 預檢請求結果的快取時間（秒）
//
// 使用示例:
//
//	config := NewCORSConfig(
//	    []string{"https://example.com"},
//	    []string{"GET", "POST"},
//	    []string{"Content-Type", "Authorization"},
//	    []string{"Content-Length"},
//	    true,
//	    3600,
//	)
//
// 簡化用法:
//
//	對於大多數應用場景，可以直接使用 DefaultCORSConfig() 獲取預設配置
func NewCORSConfig(allowOrigins, allowMethods, allowHeaders, exposeHeaders []string, allowCredentials bool, maxAge int) CORSConfig {
	return CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowMethods:     allowMethods,
		AllowHeaders:     allowHeaders,
		ExposeHeaders:    exposeHeaders,
		AllowCredentials: allowCredentials,
		MaxAge:           maxAge,
	}
}

// DefaultCORSConfig 返回預設的 CORS 配置
//
// 預設配置如下:
//   - AllowOrigins: ["*"] (允許所有來源)
//   - AllowMethods: ["GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"]
//   - AllowHeaders: ["Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"]
//   - ExposeHeaders: [""] (無)
//   - AllowCredentials: false (不允許認證信息)
//   - MaxAge: 86400 (24小時)
//
// 使用示例:
//
//	middleware := NewCORSMiddleware(DefaultCORSConfig())
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		ExposeHeaders:    []string{""},
		AllowCredentials: false,
		MaxAge:           86400,
	}
}

// CORSMiddleware 實現跨域資源共享(CORS)中間件
// 根據配置自動添加所有必要的 CORS 響應頭，處理預檢請求(OPTIONS)
type CORSMiddleware struct {
	config CORSConfig
}

// NewCORSMiddleware 創建一個新的 CORS 中間件
//
// 參數:
//   - config: CORS 配置項，可以使用 DefaultCORSConfig() 或自訂配置
//
// 使用示例:
//
//	router := gin.New()
//	corsMiddleware := NewCORSMiddleware(DefaultCORSConfig())
//	router.Use(corsMiddleware.HandlerFunc())
func NewCORSMiddleware(config CORSConfig) *CORSMiddleware {
	return &CORSMiddleware{
		config: config,
	}
}

// HandlerFunc 返回可以用於 Gin 的 CORS 中間件處理函數
//
// 此函數會:
//  1. 根據請求的 Origin 頭決定是否允許跨域請求
//  2. 添加適當的 CORS 響應頭
//  3. 自動處理預檢 OPTIONS 請求
//
// 使用示例:
//
//	router.Use(corsMiddleware.HandlerFunc())
func (m *CORSMiddleware) HandlerFunc() gin.HandlerFunc {
	// 預處理固定的頭部值以提高性能
	allowMethodsHeader := ""
	if len(m.config.AllowMethods) > 0 {
		allowMethodsHeader = headerize(m.config.AllowMethods)
	}

	allowHeadersHeader := ""
	if len(m.config.AllowHeaders) > 0 {
		allowHeadersHeader = headerize(m.config.AllowHeaders)
	}

	exposeHeadersHeader := ""
	if len(m.config.ExposeHeaders) > 0 {
		exposeHeadersHeader = headerize(m.config.ExposeHeaders)
	}

	maxAgeHeader := ""
	if m.config.MaxAge > 0 {
		maxAgeHeader = strconv.Itoa(m.config.MaxAge)
	}

	return func(c *gin.Context) {
		// 設置 Vary 頭部
		c.Writer.Header().Set("Vary", "Origin")

		// 設置允許的域名
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			for _, allowOrigin := range m.config.AllowOrigins {
				if allowOrigin == origin {
					c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
					break
				} else if allowOrigin == "*" && !m.config.AllowCredentials {
					// 只有在不允許憑證時才允許通配符
					c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
					break
				}
			}
		}

		// 設置允許的方法
		if allowMethodsHeader != "" {
			c.Writer.Header().Set("Access-Control-Allow-Methods", allowMethodsHeader)
		}

		// 設置允許的頭部信息
		if allowHeadersHeader != "" {
			c.Writer.Header().Set("Access-Control-Allow-Headers", allowHeadersHeader)
		}

		// 設置暴露的頭部信息
		if exposeHeadersHeader != "" {
			c.Writer.Header().Set("Access-Control-Expose-Headers", exposeHeadersHeader)
		}

		// 設置是否允許使用憑證
		if m.config.AllowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// 設置預檢請求的結果可以被緩存多久
		if maxAgeHeader != "" {
			c.Writer.Header().Set("Access-Control-Max-Age", maxAgeHeader)
		}

		// 處理OPTIONS請求
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// headerize 將字符串數組轉換為用逗號分隔的字符串
// 用於將多個 HTTP 頭值合併為符合 HTTP 規範的單個字符串
// 例如: ["GET", "POST"] -> "GET, POST"
func headerize(values []string) string {
	result := ""
	for i, value := range values {
		if i > 0 {
			result += ", "
		}
		result += value
	}
	return result
}
