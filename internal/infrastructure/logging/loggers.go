package logging

import (
	"go.uber.org/zap"
)

type Loggers struct {
	loggers []ILogger
}

func NewLoggers(logger ...ILogger) ILoggers {
	return &Loggers{
		loggers: logger,
	}
}

func (l *Loggers) Info(msg string, data any, fields ...zap.Field) {
	for _, logger := range l.loggers {
		logger.Info(msg, data, fields...)
	}
}
func (l *Loggers) Error(msg string, err error, data any, fields ...zap.Field) {
	for _, logger := range l.loggers {
		logger.Error(msg, err, data, fields...)
	}
}
func (l *Loggers) Debug(msg string, data any, fields ...zap.Field) {
	for _, logger := range l.loggers {
		logger.Debug(msg, data, fields...)
	}
}
func (l *Loggers) Warn(msg string, data any, fields ...zap.Field) {
	for _, logger := range l.loggers {
		logger.Warn(msg, data, fields...)
	}
}
func (l *Loggers) Write(p []byte) (n int, err error) {
	for _, logger := range l.loggers {
		logger.Debug(string(p), nil)
	}
	return len(p), nil
}

// With 設定 Identity
//
// clonedLogger := p.loggers.Clone()
//
//	newLogger := clonedLogger.With(zap.Any("identity", map[string]interface{}{
//	    "game":       p.Game,
//	    "instanceID": uuid2.New().String(),
//	}))
//
// p.loggers = newLogger
func (l *Loggers) With(fields ...zap.Field) ILoggers {
	newLoggers := make([]ILogger, 0, len(l.loggers))
	for _, logger := range l.loggers {
		// 使用 Clone 深拷貝每個 logger，並添加 identity 信息
		//newLoggers = append(newLoggers, logger.Clone().With(fields...))
		newLoggers = append(newLoggers, logger.With(fields...))
	}
	return NewLoggers(newLoggers...)
}
func (l *Loggers) Clone() ILoggers {
	clonedLoggers := make([]ILogger, len(l.loggers))
	for i, logger := range l.loggers {
		clonedLoggers[i] = logger.Clone()
	}
	return NewLoggers(clonedLoggers...)
}
