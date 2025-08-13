// Package console 提供基於 Zap 的控制台日誌輸出實現
package console

import (
	"context"
	"os"

	"github.com/tomoffice/go-clean-architecture/pkg/logger"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 定義 Logger 的配置參數
type Config struct {
	Level  logger.Level  // 日誌最低輸出等級 ("debug", "info", "warn", "error")
	Format logger.Format // 日誌輸出格式 ("json", "console")
}

// NewDefaultConfig 創建預設配置的 Logger
// 預設配置：Info 等級 + Console 格式
func NewDefaultConfig() Config {
	return Config{
		Level:  logger.InfoLevel,     // 預設 Info 等級
		Format: logger.ConsoleFormat, // 預設 Console 格式
	}
}

// Logger 實現 logger.Logger 介面，將日誌輸出到標準輸出 (stdout)
// 具體類型，可以被主套件包裝成 logger.Logger 接口
type Logger struct {
	core   zapcore.Core
	logger *zap.Logger
}

// NewLogger 根據給定配置創建新的 Logger 實例
func NewLogger(cfg Config) (*Logger, error) {
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
	level := toZapLevel(cfg.Level)
	core := zapcore.NewCore(encoder, ws, level)
	lg := zap.New(core, zap.AddCaller(), zap.AddStacktrace(toZapLevel(logger.ErrorLevel)))

	return &Logger{
		core:   core,
		logger: lg,
	}, nil
}
func NewDefaultLogger() (logger.Logger, error) {
	return NewLogger(NewDefaultConfig())
}

// Debug 輸出 Debug 等級的日誌訊息
func (l *Logger) Debug(msg string, fields ...logger.Field) {
	l.logger.Debug(msg, toZapFields(fields...)...)
}

// Info 輸出 Info 等級的日誌訊息
func (l *Logger) Info(msg string, fields ...logger.Field) {
	l.logger.Info(msg, toZapFields(fields...)...)
}

// Warn 輸出 Warn 等級的日誌訊息
func (l *Logger) Warn(msg string, fields ...logger.Field) {
	l.logger.Warn(msg, toZapFields(fields...)...)
}

// Error 輸出 Error 等級的日誌訊息
func (l *Logger) Error(msg string, fields ...logger.Field) {
	l.logger.Error(msg, toZapFields(fields...)...)
}

// With 返回一個新的 Logger 實例，預先附加指定的結構化欄位
func (l *Logger) With(fields ...logger.Field) logger.Logger {
	newZap := l.logger.With(toZapFields(fields...)...)
	return &Logger{
		core:   l.core,
		logger: newZap,
	}
}

// WithContext 從 OpenTelemetry context 中提取 trace 資訊並返回新的 Logger 實例
func (l *Logger) WithContext(ctx context.Context) logger.Logger {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		traceInfo := map[string]interface{}{
			"trace_id": spanCtx.TraceID().String(),
			"span_id":  spanCtx.SpanID().String(),
		}
		return l.With(logger.NewField("trace", traceInfo))
	}
	return l
}

// Sync 強制清空緩衝區，確保所有待輸出的日誌都已寫出
func (l *Logger) Sync() error {
	return l.logger.Sync()
}

// 確保 Logger 實現 logger.Logger 介面
var _ logger.Logger = (*Logger)(nil)
