# Gin Middleware 設計與命名規範

## 設計原則

1. **符合 Clean Architecture**：Middleware 只存在於框架層，不滲透到內層
2. **單一職責**：每個 middleware 專注於單一橫切關注點
3. **依賴抽象**：依賴內層服務的介面，而非具體實現
4. **可測試性**：每個 middleware 都可以獨立測試

## 檔案結構

```
internal/framework/http/gin/middleware/
├── middleware_interface.go  ← 定義所有共用 interface（Middleware、Logger、AuthService…）
├── auth.go                  ← JWT/API Key 驗證實作
├── cors.go                  ← CORS 實作
├── logging.go               ← Request/Response 日誌實作
├── rate_limit.go            ← Rate Limit 實作
└── recovery.go              ← Panic Recovery 實作
```

## 檔案說明

| 檔案名稱 | 功能說明 | 主要用途 |
|---------|---------|---------|
| `middleware_interface.go` | 共用介面定義 | 定義 Middleware、Logger、AuthService 等介面 |
| `auth.go` | 認證中介層 | 支援 JWT、API Key、Basic Auth 等多種認證方式 |
| `cors.go` | CORS 設定 | 處理跨域資源共享 |
| `logging.go` | 日誌記錄 | 記錄請求/回應資訊、執行時間 |
| `rate_limit.go` | 速率限制 | 限制請求頻率、防止濫用 |
| `recovery.go` | 錯誤恢復 | 處理 panic、記錄錯誤資訊 |

## 使用範例

### 認證 Middleware (JWT)

```go
// 配置 JWT 認證
authConfig := middleware.AuthConfig{
    AuthService:   authService,
    AuthType:      middleware.AuthTypeBearer,
    SecretKey:     "your-secret-key",
    TokenLookup:   "header:Authorization",
    TokenHeadName: "Bearer",
    SkipPaths: []string{
        "/api/v1/auth/login",
        "/api/v1/auth/register",
    },
}

// 使用 middleware
router := gin.New()
router.Use(middleware.NewAuthMiddleware(authConfig).Handle())

// 生成 JWT token
token, err := middleware.GenerateJWTWithExpiry(
    userID,
    username,
    email,
    roles,
    secretKey,
    24*time.Hour,
)

// 在 handler 中獲取使用者資訊
func handler(c *gin.Context) {
    userID, ok := middleware.GetUserID(c)
    if !ok {
        c.JSON(401, gin.H{"error": "未授權"})
        return
    }
    
    claims, _ := middleware.GetJWTClaims(c)
    // 使用 claims.Username, claims.Email 等
}
```

### CORS Middleware

```go
// 配置 CORS
corsConfig := middleware.CORSConfig{
    AllowOrigins:     []string{"http://localhost:3000"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
    AllowCredentials: true,
    MaxAge:           12 * time.Hour,
}

router.Use(middleware.NewCORSMiddleware(corsConfig).Handle())
```

### 日誌 Middleware

```go
// 配置日誌
logger := initLogger() // 實作 Logger 介面
router.Use(middleware.NewLoggingMiddleware(logger).Handle())
```

### 速率限制 Middleware

```go
// 配置速率限制（每分鐘 100 次請求）
router.Use(middleware.NewRateLimitMiddleware(100, time.Minute, nil).Handle())
```

### Recovery Middleware

```go
// 配置 panic 恢復
recoveryConfig := middleware.RecoveryConfig{
    Logger:     logger,
    Stack:      true,
    PrintStack: gin.IsDebugging(),
}
router.Use(middleware.NewRecoveryMiddleware(recoveryConfig).Handle())
```

## 在專案中整合所有 Middleware

```go
func main() {
    // 初始化服務
    logger := initLogger()
    authService := initAuthService()
    
    // 建立路由器
    router := gin.New()
    
    // 應用 middleware（按照執行順序）
    router.Use(middleware.NewRecoveryMiddleware(recoveryConfig).Handle())
    router.Use(middleware.NewLoggingMiddleware(logger).Handle())
    router.Use(middleware.NewCORSMiddleware(corsConfig).Handle())
    router.Use(middleware.NewRateLimitMiddleware(100, time.Minute, nil).Handle())
    router.Use(middleware.NewAuthMiddleware(authConfig).Handle())
    
    // 設置路由
    setupRoutes(router)
    
    // 啟動服務器
    router.Run(":8080")
}
```

## 開發新的 Middleware

建立新的 middleware 時，請遵循以下步驟：

1. 在 `middleware_interface.go` 中定義任何需要的新介面
2. 創建新的檔案，如 `new_feature.go`
3. 實現 `Middleware` 介面
4. 提供建構函數 `NewXxxMiddleware`
5. 撰寫單元測試
6. 更新此文件

```go
type NewFeatureMiddleware struct {
    // 依賴的服務
}

func NewNewFeatureMiddleware(/* dependencies */) Middleware {
    return &NewFeatureMiddleware{
        // 初始化
    }
}

func (m *NewFeatureMiddleware) Handle() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 前置處理
        
        c.Next() // 執行下一個 handler
        
        // 後置處理（如果需要）
    }
}
```

## 測試指南

每個 middleware 都應該有對應的測試檔案：
- `auth_test.go`
- `cors_test.go`
- `logging_test.go`
- `rate_limit_test.go`
- `recovery_test.go`

測試應涵蓋：
- 正常情況
- 錯誤處理
- 邊界條件
- 效能測試（對於 rate limit 等）
