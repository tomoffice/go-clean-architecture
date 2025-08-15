# pkg/logger

## 用途與職責

高度可擴展的結構化日誌套件，為整個專案提供統一的日誌記錄介面。支援多種輸出 adapter、OpenTelemetry 整合，以及 Tee 模式同時輸出到多個目標，確保日誌記錄的一致性和可觀測性。

## 設計原則

- **介面統一**: 提供統一的 Logger 介面，隔離具體實現
- **多 Adapter 支援**: Console、GCP Cloud Logging、Seq 等多種輸出方式
- **結構化日誌**: 基於 zap 的高效能結構化日誌記錄
- **OpenTelemetry 整合**: 自動提取和記錄分散式追蹤資訊
- **Tee 模式**: 支援同時輸出到多個目標
- **型別安全**: 完整的 Go 型別系統支援

## 目錄結構

```
pkg/logger/
├── logger.go                    # Logger 介面定義
├── field.go                     # Field 工具函數
├── level.go                     # 日誌等級定義和解析
├── tee.go                       # Tee logger 實作
├── mock/                        # Mock 實作
│   └── mock_logger.go
├── adapters/                    # Logger adapters
│   ├── console/                 # Console 輸出
│   ├── gcp/                     # Google Cloud Logging
│   └── seq/                     # Seq 日誌平台
└── tests/                       # 整合測試
    └── integration/
```

## 使用方式

### 基本使用

```go
import (
    "github.com/tomoffice/go-clean-architecture/pkg/logger"
    "github.com/tomoffice/go-clean-architecture/pkg/logger/adapters/console"
)

// 創建 logger
lg, err := console.NewDefaultLogger()
if err != nil {
    panic(err)
}
defer lg.Sync()

// 結構化日誌記錄
lg.Info("Hello World", logger.NewField("user", "alice"))
lg.WithContext(ctx).Info("Processing request")  // 自動包含 trace 資訊
```

### 多 Adapter 組合

```go
// 創建多個 logger 並組合
consoleLogger, _ := console.NewDefaultLogger()
seqLogger, _ := seq.NewLogger(seq.Config{Endpoint: "http://localhost:5341"})

// Tee 模式同時輸出
combinedLogger := logger.NewTeeLogger(consoleLogger, seqLogger)
combinedLogger.Info("Message sent to both outputs")
```

### 依賴注入使用

```go
// 在 UseCase 中接收 Logger 介面
type MemberUseCase struct {
    logger logger.Logger
}

func (uc *MemberUseCase) CreateMember(ctx context.Context, req dto.CreateMemberDTO) error {
    contextLogger := uc.logger.WithContext(ctx)  // 包含 trace 資訊
    contextLogger.Info("Creating member", logger.NewField("email", req.Email))
    
    // 業務邏輯...
    
    contextLogger.Info("Member created successfully", logger.NewField("member_id", memberID))
    return nil
}
```

## 注意事項

- **統一介面**: 所有層級都應使用 logger.Logger 介面，不直接依賴具體實現
- **Context 傳遞**: 使用 WithContext() 確保分散式追蹤資訊正確記錄
- **結構化欄位**: 優先使用 logger.NewField() 而非字串插值
- **資源管理**: 在程式結束前必須呼叫 Sync() 確保日誌完整輸出
- **測試支援**: 提供完整的 Mock 實現供測試使用