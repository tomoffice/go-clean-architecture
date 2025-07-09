package gcp

import (
	"cloud.google.com/go/logging"
	"go.uber.org/zap/zapcore"
)

func mapSeverity(l zapcore.Level) logging.Severity {
	switch l {
	case zapcore.DebugLevel:
		return logging.Debug
	case zapcore.InfoLevel:
		return logging.Info
	case zapcore.WarnLevel:
		return logging.Warning
	case zapcore.ErrorLevel:
		return logging.Error
	case zapcore.DPanicLevel, zapcore.PanicLevel:
		return logging.Critical
	case zapcore.FatalLevel:
		return logging.Critical
	default:
		return logging.Default
	}
}
