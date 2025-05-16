package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ILoggerBase interface {
	BInfo(msg string, fields ...zap.Field)
	BError(msg string, fields ...zap.Field)
	BDebug(msg string, fields ...zap.Field)
	BWarn(msg string, fields ...zap.Field)
	IsLevelEnabled(level zapcore.Level) bool
}
type ILogger interface {
	Info(msg string, data any, fields ...zap.Field)
	Error(msg string, err error, data any, fields ...zap.Field)
	Debug(msg string, data any, fields ...zap.Field)
	Warn(msg string, data any, fields ...zap.Field)
	Write(p []byte) (n int, err error)
	With(fields ...zap.Field) ILogger
	Clone() ILogger
}

type ILoggers interface {
	Info(msg string, data any, fields ...zap.Field)
	Error(msg string, err error, data any, fields ...zap.Field)
	Debug(msg string, data any, fields ...zap.Field)
	Warn(msg string, data any, fields ...zap.Field)
	Write(p []byte) (n int, err error)
	With(fields ...zap.Field) ILoggers
	Clone() ILoggers
}
