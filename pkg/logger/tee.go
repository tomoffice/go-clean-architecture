package logger

import (
	"context"
	"sync"
)

// TeeLogger 實現將日誌同時輸出到多個 Logger 的功能
// 類似 Unix 的 tee 命令，將輸入分流到多個輸出
type TeeLogger struct {
	loggers []Logger
	mu      sync.RWMutex
}

// 確保 teeLogger 實現 Logger 介面
var _ Logger = (*TeeLogger)(nil)

// NewTeeLogger 創建新的 tee logger，將日誌分流到多個 logger
// 公開函數，允許應用層組合多個 logger
func NewTeeLogger(loggers ...Logger) *TeeLogger {
	if len(loggers) == 0 {
		// 返回空的 tee logger，所有操作都是 no-op
		return &TeeLogger{loggers: []Logger{}}
	}

	// 創建副本以避免外部修改
	copied := make([]Logger, len(loggers))
	copy(copied, loggers)

	return &TeeLogger{
		loggers: copied,
	}
}

// Debug 將 Debug 日誌分流到所有 logger
func (t *TeeLogger) Debug(msg string, fields ...Field) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, lg := range t.loggers {
		lg.Debug(msg, fields...)
	}
}

// Info 將 Info 日誌分流到所有 logger
func (t *TeeLogger) Info(msg string, fields ...Field) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, lg := range t.loggers {
		lg.Info(msg, fields...)
	}
}

// Warn 將 Warn 日誌分流到所有 logger
func (t *TeeLogger) Warn(msg string, fields ...Field) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, lg := range t.loggers {
		lg.Warn(msg, fields...)
	}
}

// Error 將 Error 日誌分流到所有 logger
func (t *TeeLogger) Error(msg string, fields ...Field) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	for _, lg := range t.loggers {
		lg.Error(msg, fields...)
	}
}

// With 對所有 logger 應用 With，返回新的 teeLogger
func (t *TeeLogger) With(fields ...Field) Logger {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// 對每個 logger 應用 With，創建新的 tee logger
	newLoggers := make([]Logger, len(t.loggers))
	for i, lg := range t.loggers {
		newLoggers[i] = lg.With(fields...)
	}

	return &TeeLogger{
		loggers: newLoggers,
	}
}

// WithContext 對所有 logger 應用 WithContext，返回新的 teeLogger
func (t *TeeLogger) WithContext(ctx context.Context) Logger {
	t.mu.RLock()
	defer t.mu.RUnlock()

	// 對每個 logger 應用 WithContext，創建新的 tee logger
	newLoggers := make([]Logger, len(t.loggers))
	for i, lg := range t.loggers {
		newLoggers[i] = lg.WithContext(ctx)
	}

	return &TeeLogger{
		loggers: newLoggers,
	}
}

// Sync 對所有 logger 執行 Sync，返回最後一個錯誤
func (t *TeeLogger) Sync() error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var lastErr error
	for _, lg := range t.loggers {
		if err := lg.Sync(); err != nil {
			lastErr = err
		}
	}

	return lastErr
}
