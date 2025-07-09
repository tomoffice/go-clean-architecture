// Package console 提供基於 Zap 的控制台日誌輸出實現
package console

import (
	"os"

	"github.com/tomoffice/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 定義 ConsoleLogger 的配置參數
type Config struct {
	Level  zapcore.Level // 日誌最低輸出等級 (DebugLevel, InfoLevel, WarnLevel, ErrorLevel)
	Format logger.Format // 日誌輸出格式，支援 JSONFormat 或 ConsoleFormat
}

// NewDefaultConfig 創建預設配置的 ConsoleLogger
// 預設配置：Info 等級 + Console 格式
func NewDefaultConfig() (logger.Logger, error) {
	return NewLogger(Config{
		Level:  zapcore.InfoLevel,
		Format: logger.ConsoleFormat,
	})
}

// Logger 實現 logger.Logger 介面，將日誌輸出到標準輸出 (stdout)
type Logger struct {
	core   zapcore.Core
	logger *zap.Logger
}

// NewLogger 根據給定配置創建新的 ConsoleLogger 實例
func NewLogger(cfg Config) (logger.Logger, error) {
	// 1. 建立統一的編碼器配置
	encCfg := zapcore.EncoderConfig{
		TimeKey:       "timestamp",
		LevelKey:      "level",
		CallerKey:     "caller",
		MessageKey:    "message",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}

	// 2. 根據格式選擇對應的編碼器
	var encoder zapcore.Encoder
	switch cfg.Format {
	case logger.JSONFormat:
		encoder = zapcore.NewJSONEncoder(encCfg)
	case logger.ConsoleFormat:
		encoder = zapcore.NewConsoleEncoder(encCfg)
	default:
		return nil, logger.ErrUnsupportedFormat
	}

	// 3. 建立 Core 和 Logger 實例
	ws := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(encoder, ws, cfg.Level)
	lg := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{
		core:   core,
		logger: lg,
	}, nil
}

// Debug 輸出 Debug 等級的日誌訊息
func (l *Logger) Debug(msg string, fields ...logger.Field) {
	l.logger.Debug(msg, fields...)
}

// Info 輸出 Info 等級的日誌訊息
func (l *Logger) Info(msg string, fields ...logger.Field) {
	l.logger.Info(msg, fields...)
}

// Warn 輸出 Warn 等級的日誌訊息
func (l *Logger) Warn(msg string, fields ...logger.Field) {
	l.logger.Warn(msg, fields...)
}

// Error 輸出 Error 等級的日誌訊息
func (l *Logger) Error(msg string, fields ...logger.Field) {
	l.logger.Error(msg, fields...)
}

// With 返回一個新的 Logger 實例，預先附加指定的結構化欄位
func (l *Logger) With(fields ...logger.Field) logger.Logger {
	newZap := l.logger.With(fields...)
	return &Logger{
		core:   l.core,
		logger: newZap,
	}
}

// Sync 強制清空緩衝區，確保所有待輸出的日誌都已寫出
func (l *Logger) Sync() error {
	return l.logger.Sync()
}

// GetCore 回傳底層 zapcore.Core，供 MultiLogger 組合使用
func (l *Logger) GetCore() zapcore.Core {
	return l.core
}
