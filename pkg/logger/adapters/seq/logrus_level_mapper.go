package seq

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/zap/zapcore"
)

// mapLevel 將 zapcore.Level 轉成 logrus.Level
func mapLevel(zl zapcore.Level) logrus.Level {
	switch zl {
	case zapcore.DebugLevel:
		return logrus.DebugLevel
	case zapcore.InfoLevel:
		return logrus.InfoLevel
	case zapcore.WarnLevel:
		return logrus.WarnLevel
	case zapcore.ErrorLevel:
		return logrus.ErrorLevel
	case zapcore.DPanicLevel, zapcore.PanicLevel:
		return logrus.PanicLevel
	case zapcore.FatalLevel:
		return logrus.FatalLevel
	default:
		return logrus.InfoLevel
	}
}