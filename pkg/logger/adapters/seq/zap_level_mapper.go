package seq

import (
	"go.uber.org/zap/zapcore"

	"github.com/tomoffice/go-clean-architecture/pkg/logger"
)

// toZapLevel 將抽象的 logger.Level 轉換為 zapcore.Level
// 這確保了抽象層不洩漏 zapcore 實作細節
func toZapLevel(level logger.Level) zapcore.Level {
	switch level {
	case logger.DebugLevel:
		return zapcore.DebugLevel
	case logger.InfoLevel:
		return zapcore.InfoLevel
	case logger.WarnLevel:
		return zapcore.WarnLevel
	case logger.ErrorLevel:
		return zapcore.ErrorLevel
	case logger.PanicLevel:
		return zapcore.PanicLevel
	case logger.FatalLevel:
		return zapcore.FatalLevel
	default:
		// 預設返回 Info 等級
		return zapcore.InfoLevel
	}
}