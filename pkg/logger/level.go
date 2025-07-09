package logger

// 日誌等級常數定義（InfoLevel 等）

import "go.uber.org/zap/zapcore"

type Level = zapcore.Level

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	PanicLevel = zapcore.PanicLevel
	FatalLevel = zapcore.FatalLevel
)

// ParseLevel 解析字串為 Level
func ParseLevel(s string) (Level, error) {
	return zapcore.ParseLevel(s)
}
