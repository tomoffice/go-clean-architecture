package logger

import "go.uber.org/zap/zapcore"

// Level 代表日誌的重要性等級。
// 等級由低到高依序為：Debug < Info < Warn < Error < Panic < Fatal
//
// 日誌等級的用途：
//   - Debug: 詳細的除錯資訊，通常只在開發環境啟用
//   - Info:  一般的應用程式資訊，記錄正常的操作流程
//   - Warn:  警告訊息，表示潛在問題但不影響程式繼續執行
//   - Error: 錯誤訊息，表示程式遇到錯誤但可以恢復
//   - Panic: 嚴重錯誤，會引發 panic 中斷程式執行
//   - Fatal: 致命錯誤，會終止整個程式
//
// 當設定日誌等級時，只有等於或高於設定等級的日誌才會被輸出。
// 例如設定為 InfoLevel 時，Debug 等級的日誌不會被輸出。
type Level = zapcore.Level

// 日誌等級常數定義
const (
	// DebugLevel 除錯等級，用於詳細的程式運行資訊
	DebugLevel = zapcore.DebugLevel

	// InfoLevel 資訊等級，用於一般的程式狀態記錄
	InfoLevel = zapcore.InfoLevel

	// WarnLevel 警告等級，用於潛在問題的提醒
	WarnLevel = zapcore.WarnLevel

	// ErrorLevel 錯誤等級，用於程式錯誤的記錄
	ErrorLevel = zapcore.ErrorLevel

	// PanicLevel Panic 等級，記錄後會引發 panic
	PanicLevel = zapcore.PanicLevel

	// FatalLevel 致命等級，記錄後會終止程式執行
	FatalLevel = zapcore.FatalLevel
)

// ParseLevel 將日誌等級字串解析為對應的 Level 常數。
// 此函數不區分大小寫，並提供容錯處理。
//
// 支援的字串值：
//   - "debug" -> DebugLevel
//   - "info"  -> InfoLevel
//   - "warn"  -> WarnLevel
//   - "error" -> ErrorLevel
//   - "panic" -> PanicLevel
//   - "fatal" -> FatalLevel
//
// 參數：
//   - level: 日誌等級的字串表示
//
// 返回值：
//   - Level: 對應的日誌等級常數，無效輸入時返回 InfoLevel
//
// 使用範例：
//
//	level := logger.ParseLevel("debug")  // 返回 DebugLevel
//	level := logger.ParseLevel("invalid") // 返回 InfoLevel (預設值)
//
//	// 在配置中使用
//	cfg := console.Config{
//	    Level: logger.ParseLevel(os.Getenv("LOG_LEVEL")),
//	}
func ParseLevel(level string) Level {
	switch level {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	case "panic":
		return PanicLevel
	case "fatal":
		return FatalLevel
	default:
		// 無效輸入時返回 Info 等級作為安全預設值
		return InfoLevel
	}
}
