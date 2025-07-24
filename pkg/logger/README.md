# Logger 套件

高度可擴展的結構化日誌套件，支援多種輸出 adapter 和 OpenTelemetry 整合。

## 特性

- **多 Adapter 支援**: Console、GCP Cloud Logging、Seq
- **結構化日誌**: 基於 zap 的高效能日誌記錄
- **OpenTelemetry 整合**: 自動提取和記錄 trace 資訊
- **Tee 模式**: 同時輸出到多個目標
- **型別安全**: 完整的 Go 型別系統支援
- **可配置**: 支援多種日誌等級和格式

## 快速開始

### 基本使用

```go
package main

import (
    "github.com/tomoffice/go-clean-architecture/pkg/logger"
    "github.com/tomoffice/go-clean-architecture/pkg/logger/adapters/console"
)

func main() {
    // 創建 console logger
    lg, err := console.NewLogger(console.Config{
        Level:  logger.InfoLevel,
        Format: logger.ConsoleFormat,
    })
    if err != nil {
        panic(err)
    }
    defer lg.Sync()

    // 基本日誌記錄
    lg.Info("Hello World", logger.NewField("user", "alice"))
    lg.Warn("Warning message", logger.NewField("code", 404))
    lg.Error("Error occurred", logger.NewField("error", "connection failed"))
}
```

### 多 Adapter 組合

```go
package main

import (
    "github.com/tomoffice/go-clean-architecture/pkg/logger"
    "github.com/tomoffice/go-clean-architecture/pkg/logger/adapters/console"
    "github.com/tomoffice/go-clean-architecture/pkg/logger/adapters/seq"
)

func main() {
    // 創建多個 logger
    consoleLogger, _ := console.NewLogger(console.Config{
        Level:  logger.InfoLevel,
        Format: logger.JSONFormat,
    })
    
    seqLogger, _ := seq.NewLogger(seq.Config{
        Endpoint: "http://localhost:5341",
        Level:    logger.InfoLevel,
    })

    // 使用 tee 組合多個 logger
    combinedLogger := logger.NewTeeLogger(consoleLogger, seqLogger)
    
    // 日誌會同時輸出到 console 和 Seq
    combinedLogger.Info("Message sent to both outputs")
}
```

### OpenTelemetry 整合

```go
package main

import (
    "context"
    "go.opentelemetry.io/otel"
    "github.com/tomoffice/go-clean-architecture/pkg/logger"
    "github.com/tomoffice/go-clean-architecture/pkg/logger/adapters/console"
)

func main() {
    lg, _ := console.NewDefaultLogger()
    tracer := otel.Tracer("my-service")

    // 創建 span
    ctx, span := tracer.Start(context.Background(), "main-operation")
    defer span.End()

    // WithContext 會自動提取 trace 資訊
    lg.WithContext(ctx).Info("Operation started")
    // 輸出: {"level":"info","message":"Operation started","trace":{"trace_id":"...","span_id":"..."}}
}
```

## 架構設計

### 核心介面

```go
// Logger 主要介面
type Logger interface {
    Debug(msg string, fields ...Field)
    Info(msg string, fields ...Field)
    Warn(msg string, fields ...Field)
    Error(msg string, fields ...Field)
    With(fields ...Field) Logger
    WithContext(ctx context.Context) Logger
    Sync() error
}

// 終止程序的擴展介面
type TerminatingLogger interface {
    Logger
    Panic(msg string, fields ...Field)
    Fatal(msg string, fields ...Field)
}

// 可選的 Core 存取介面
type Core interface {
    GetCore() zapcore.Core
}
```

### 套件結構

```
pkg/logger/
├── logger.go                    # Logger 介面定義
├── field.go                     # Field 工具函數
├── level.go                     # 日誌等級定義和解析
├── format.go                    # 日誌格式定義和解析
├── errors.go                    # 錯誤定義
├── config.go                    # 配置結構
├── tee.go                       # Tee logger 實作
├── mock/                        # Mock 實作
│   └── mock_logger.go
├── adapters/                    # Logger adapters
│   ├── console/                 # Console 輸出
│   │   └── console_logger.go
│   ├── gcp/                     # Google Cloud Logging
│   │   ├── gcp_logger.go
│   │   └── serverity_mapper.go
│   └── seq/                     # Seq 日誌平台
│       ├── seq_logger.go
│       ├── seq_core.go
│       ├── seq_sender.go
│       ├── sender.go
│       ├── field_helper.go
│       ├── level_mapper.go
│       └── mock/
└── tests/                       # 整合測試
    └── integration/
        └── seq_integration_test.go
```

## Adapters

### Console Adapter

輸出到標準輸出，支援 JSON 和 Console 格式。

```go
import "github.com/tomoffice/go-clean-architecture/pkg/logger/adapters/console"

// 創建 console logger
logger, err := console.NewLogger(console.Config{
    Level:  logger.InfoLevel,      // debug, info, warn, error
    Format: logger.ConsoleFormat,  // console, json
})

// 或使用預設配置
logger, err := console.NewDefaultLogger()
```

### GCP Cloud Logging Adapter

整合 Google Cloud Logging 服務。

```go
import "github.com/tomoffice/go-clean-architecture/pkg/logger/adapters/gcp"

logger, err := gcp.NewLogger(gcp.Config{
    ProjectID: "my-gcp-project",
    LogName:   "my-app-logs",
    Level:     logger.InfoLevel,
})
```

### Seq Adapter

整合 Seq 結構化日誌平台。

```go
import "github.com/tomoffice/go-clean-architecture/pkg/logger/adapters/seq"

logger, err := seq.NewLogger(seq.Config{
    Endpoint: "http://localhost:5341",
    APIKey:   "your-api-key",        // 可選
    Level:    logger.InfoLevel,
})
```

## API 參考

### 日誌等級

```go
const (
    DebugLevel = zapcore.DebugLevel
    InfoLevel  = zapcore.InfoLevel
    WarnLevel  = zapcore.WarnLevel
    ErrorLevel = zapcore.ErrorLevel
    PanicLevel = zapcore.PanicLevel
    FatalLevel = zapcore.FatalLevel
)

// 從字串解析等級
level := logger.ParseLevel("info")  // 返回 logger.InfoLevel
```

### 日誌格式

```go
const (
    JSONFormat    Format = "json"
    ConsoleFormat Format = "console"
)

// 從字串解析格式
format := logger.ParseFormat("json")  // 返回 logger.JSONFormat
```

### Field 建立

```go
// 通用 field 建立函數
field := logger.NewField("key", "value")
field := logger.NewField("count", 42)
field := logger.NewField("enabled", true)
field := logger.NewField("user", map[string]string{"name": "alice"})
```

### 錯誤定義

```go
var (
    ErrNoValidLoggers        = errors.New("沒有有效的 logger 可用")
    ErrInvalidConfig         = errors.New("無效的 logger 設定")
    ErrLoggerCreationFailed  = errors.New("建立 logger 失敗")
    ErrUnsupportedFormat     = errors.New("不支援的日誌格式")
)
```

## OpenTelemetry 整合

### Trace 資訊自動提取

當使用 `WithContext()` 方法時，logger 會自動從 OpenTelemetry context 中提取 trace 資訊：

```go
// 日誌會包含結構化的 trace 資訊
lg.WithContext(ctx).Info("Processing request")

// 輸出格式:
{
  "level": "info",
  "message": "Processing request",
  "trace": {
    "trace_id": "4bf92f3577b34da6a3ce929d0e0e4736",
    "span_id": "00f067aa0ba902b7"
  },
  "timestamp": "2024-01-01T12:00:00Z"
}
```

### Distributed Tracing 支援

```go
func processOrder(ctx context.Context, orderID string) {
    tracer := otel.Tracer("order-service")
    
    // 為每個操作創建子 span
    ctx, span := tracer.Start(ctx, "validate-order")
    logger.WithContext(ctx).Info("Validating order", 
        logger.NewField("order_id", orderID))
    span.End()
    
    ctx, span = tracer.Start(ctx, "process-payment")
    logger.WithContext(ctx).Info("Processing payment",
        logger.NewField("order_id", orderID))
    span.End()
}
```

## 測試

### 執行測試

```bash
# 執行所有測試
go test ./...

# 執行整合測試 (需要 Seq 服務運行)
go test -tags=integration ./tests/integration/

# 執行特定測試
go test ./adapters/console/
```

### 使用 Mock

```go
import "github.com/tomoffice/go-clean-architecture/pkg/logger/mock"

// 在測試中使用 mock logger
mockLogger := mock.NewMockLogger(ctrl)
mockLogger.EXPECT().Info("test message", gomock.Any()).Times(1)

// 測試你的代碼
yourFunction(mockLogger)
```

### 整合測試範例

```go
func TestSeqIntegration(t *testing.T) {
    // 需要 Seq 服務在 localhost:5341 運行
    cfg := seq.Config{
        Endpoint: "http://localhost:5341",
        Level:    logger.InfoLevel,
    }
    
    lg, err := seq.NewLogger(cfg)
    require.NoError(t, err)
    defer lg.Sync()
    
    lg.Info("Test message", logger.NewField("test", true))
}
```

## 最佳實踐

### 1. 錯誤處理

```go
// 總是檢查 logger 建立錯誤
logger, err := console.NewLogger(config)
if err != nil {
    log.Fatalf("Failed to create logger: %v", err)
}

// 確保在程式結束前同步
defer logger.Sync()
```

### 2. 結構化欄位

```go
// 使用結構化欄位而非字串插值
// 好的做法
logger.Info("User login", 
    logger.NewField("user_id", userID),
    logger.NewField("ip", clientIP))

// 避免的做法
logger.Info(fmt.Sprintf("User %s login from %s", userID, clientIP))
```

### 3. Context 傳遞

```go
// 在整個請求生命週期中傳遞 context
func handleRequest(ctx context.Context, logger logger.Logger) {
    // 使用 WithContext 確保 trace 資訊被記錄
    contextLogger := logger.WithContext(ctx)
    
    contextLogger.Info("Request started")
    // ... 處理邏輯
    contextLogger.Info("Request completed")
}
```

### 4. 效能考量

```go
// 重用 logger 實例
var globalLogger logger.Logger

// 使用 With() 創建帶有共同欄位的 logger
userLogger := globalLogger.With(
    logger.NewField("user_id", userID),
    logger.NewField("session_id", sessionID),
)

// 在整個用戶會話中重用 userLogger
userLogger.Info("Action performed")
```