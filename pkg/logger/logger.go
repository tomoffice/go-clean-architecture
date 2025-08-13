//go:generate mockgen -source=logger.go -destination=mock/mock_logger.go -package=mock

// Package logger 提供純抽象的結構化日誌介面，不依賴任何第三方實作。
//
// 主要特性：
//   - 純抽象介面：不洩漏底層實作細節
//   - 結構化日誌：支援結構化欄位附加
//   - 可觀測性整合：支援從 context 提取觀測欄位
//   - 可替換性：底層實作完全可替換
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

import "context"

// Logger 定義結構化日誌記錄的抽象介面。
// 這是套件的核心介面，所有 adapter 都必須實現此介面以提供一致的日誌記錄行為。
//
// Logger 支援四種標準日誌等級（Debug、Info、Warn、Error），
// 以及結構化欄位附加和可觀測性 context 整合。
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

	// WithContext 從 context 中提取與此實作相關的觀測欄位並返回新的 Logger 實例。
	// 具體提取的欄位內容由底層實作決定，可能包含 trace id、span id 或其他觀測資訊。
	//
	// 此方法保持抽象層的中性，不承諾特定的觀測框架或欄位格式。
	WithContext(ctx context.Context) Logger

	// Sync 強制清空所有緩衝的日誌到底層輸出。
	// 建議在程式結束前呼叫此方法以確保所有日誌都被正確寫出。
	Sync() error
}

// TerminatingLogger 擴展 Logger 介面，新增可終止程式執行的日誌方法。
// 
// 注意：這些方法通常應該在應用程式的最外層使用（如 main 函數），
// 業務邏輯層一般不應該直接決定程式的終止行為。
type TerminatingLogger interface {
	Logger

	// Panic 記錄 Panic 等級的日誌訊息並引發 panic。
	// 此方法會在記錄日誌後呼叫 panic()，用於嚴重的程式錯誤。
	Panic(msg string, fields ...Field)

	// Fatal 記錄 Fatal 等級的日誌訊息並終止程式。
	// 此方法會在記錄日誌後呼叫 os.Exit(1)，用於無法恢復的致命錯誤。
	Fatal(msg string, fields ...Field)
}