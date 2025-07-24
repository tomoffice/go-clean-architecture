//go:generate mockgen -source=logger.go -destination=mock/mock_logger.go -package=mock

// Package logger 提供高度可擴展的結構化日誌套件，支援多種輸出 adapter 和 OpenTelemetry 整合。
//
// 主要特性：
//   - 多 Adapter 支援：Console、GCP Cloud Logging、Seq
//   - 結構化日誌：基於 zap 的高效能日誌記錄
//   - OpenTelemetry 整合：自動提取和記錄 trace 資訊
//   - Tee 模式：同時輸出到多個目標
//   - 型別安全：完整的 Go 型別系統支援
//
// 基本使用範例：
//
//	import "github.com/tomoffice/go-clean-architecture/pkg/logger"
//	import "github.com/tomoffice/go-clean-architecture/pkg/logger/adapters/console"
//
//	lg, err := console.NewLogger(console.Config{
//	    Level:  logger.InfoLevel,
//	    Format: logger.ConsoleFormat,
//	})
//	if err != nil {
//	    panic(err)
//	}
//	defer lg.Sync()
//
//	lg.Info("Hello World", logger.NewField("user", "alice"))
package logger

import (
	"context"
	"go.uber.org/zap/zapcore"
)

// Logger 定義結構化日誌記錄的抽象介面。
// 這是套件的核心介面，所有 adapter 都必須實現此介面以提供一致的日誌記錄行為。
//
// Logger 支援四種標準日誌等級（Debug、Info、Warn、Error），
// 以及結構化欄位附加和 OpenTelemetry context 整合。
type Logger interface {
	// Debug 記錄除錯等級的日誌訊息，通常用於開發階段的詳細資訊。
	Debug(msg string, fields ...Field)

	// Info 記錄資訊等級的日誌訊息，用於一般的應用程式狀態和操作記錄。
	Info(msg string, fields ...Field)

	// Warn 記錄警告等級的日誌訊息，用於可能影響程式運行但不會導致錯誤的情況。
	Warn(msg string, fields ...Field)

	// Error 記錄錯誤等級的日誌訊息，用於程式錯誤和異常情況。
	Error(msg string, fields ...Field)

	// With 返回一個新的 Logger 實例，預先附加指定的結構化欄位。
	// 這個方法不會修改原始的 Logger，而是返回一個新的實例。
	//
	// 範例：
	//   userLogger := logger.With(logger.NewField("user_id", "123"))
	//   userLogger.Info("User action") // 自動包含 user_id 欄位
	With(fields ...Field) Logger

	// WithContext 從 OpenTelemetry context 中提取 trace 資訊並返回新的 Logger 實例。
	// 如果 context 包含有效的 span，trace_id 和 span_id 會被自動添加到日誌中。
	//
	// 輸出格式：
	//   {"level":"info","message":"Processing","trace":{"trace_id":"...","span_id":"..."}}
	WithContext(ctx context.Context) Logger

	// Sync 強制清空所有緩衝的日誌到底層輸出。
	// 建議在程式結束前呼叫此方法以確保所有日誌都被正確寫出。
	Sync() error
}

// TerminatingLogger 擴展 Logger 介面，新增可終止程式執行的日誌方法。
// 實作此介面的 Logger 可以在記錄日誌後終止程式執行。
type TerminatingLogger interface {
	Logger

	// Panic 記錄 Panic 等級的日誌訊息並引發 panic。
	// 此方法會在記錄日誌後呼叫 panic()，用於嚴重的程式錯誤。
	Panic(msg string, fields ...Field)

	// Fatal 記錄 Fatal 等級的日誌訊息並終止程式。
	// 此方法會在記錄日誌後呼叫 os.Exit(1)，用於無法恢復的致命錯誤。
	Fatal(msg string, fields ...Field)
}

// Core 是可選介面，用於需要存取底層 zapcore.Core 的進階使用場景。
// 只有需要直接操作底層 Core 的客戶端才應該依賴此介面，
// 一般使用情況下不需要使用此介面。
//
// 注意：使用此介面會增加對 zap 的依賴，可能影響套件的可移植性。
type Core interface {
	// GetCore 返回底層的 zapcore.Core 實例。
	// 此方法主要用於需要組合多個 Core 或進行底層操作的場景。
	GetCore() zapcore.Core
}
