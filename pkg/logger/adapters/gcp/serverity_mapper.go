package gcp

import (
	"cloud.google.com/go/logging"
	"github.com/tomoffice/go-clean-architecture/pkg/logger"
	"go.uber.org/zap/zapcore"
)

func mapSeverity(l logger.Level) logging.Severity {
	switch zapcore.Level(l) {
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
