// Package seq 提供基於 Zap 和 Logrus 的 Seq 日誌服務實現
// 透過 Bridge 模式將 Zap 日誌轉發到 Seq 日誌平台
package seq

import (
	"context"
	"github.com/tomoffice/go-clean-architecture/pkg/logger"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 定義 Seq Logger 的配置參數
type Config struct {
	Endpoint string       // Seq 服務端點，例如：http://seq:5341
	APIKey   string       // Seq API 金鑰，可選參數
	Level    logger.Level // 日誌最低輸出等級
}

// NewDefaultConfig 創建預設配置的 Seq Logger
// 預設配置：本地端點 localhost:5341 + Info 等級
func NewDefaultConfig() Config {
	return Config{
		Endpoint: "http://localhost:5341",
		APIKey:   "",
		Level:    logger.InfoLevel,
	}
}

// Logger 實現 logger.Logger 介面，提供 Seq 日誌輸出功能
type Logger struct {
	core   zapcore.Core
	logger *zap.Logger
}

// NewLogger 創建新的 Seq Logger 實例
func NewLogger(cfg Config) (logger.Logger, error) {
	// 建立使用Logrus 的 Seq Sender
	sender := NewLogrusSender(cfg.Endpoint, cfg.APIKey)

	encCfg := zapcore.EncoderConfig{
		TimeKey:        "@t",  // Seq 時間戳欄位
		LevelKey:       "@l",  // Seq 等級欄位
		MessageKey:     "@mt", // Seq 訊息範本欄位
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout(time.RFC3339Nano),
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}

	enc := zapcore.NewJSONEncoder(encCfg)

	core := NewSeqCore(sender, enc, cfg.Level)
	// 4) 建 zap.Logger
	zl := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{
		core:   core,
		logger: zl,
	}, nil
}
func NewDefaultLogger() (logger.Logger, error) {
	return NewLogger(NewDefaultConfig())
}

// Debug 輸出 Debug 等級的日誌訊息
func (s *Logger) Debug(msg string, fields ...logger.Field) { s.logger.Debug(msg, fields...) }

// Info 輸出 Info 等級的日誌訊息
func (s *Logger) Info(msg string, fields ...logger.Field) { s.logger.Info(msg, fields...) }

// Warn 輸出 Warn 等級的日誌訊息
func (s *Logger) Warn(msg string, fields ...logger.Field) { s.logger.Warn(msg, fields...) }

// Error 輸出 Error 等級的日誌訊息
func (s *Logger) Error(msg string, fields ...logger.Field) { s.logger.Error(msg, fields...) }

// With 返回一個新的 Logger 實例，預先附加指定的結構化欄位
func (s *Logger) With(fields ...logger.Field) logger.Logger {
	newZap := s.logger.With(fields...)
	return &Logger{
		core:   s.core, // core 保持不變
		logger: newZap, // 只更新 logger
	}
}

// WithContext 從 OpenTelemetry context 中提取 trace 資訊並返回新的 Logger 實例
func (s *Logger) WithContext(ctx context.Context) logger.Logger {
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.IsValid() {
		traceInfo := map[string]interface{}{
			"trace_id": spanCtx.TraceID().String(),
			"span_id":  spanCtx.SpanID().String(),
		}
		return s.With(logger.NewField("trace", traceInfo))
	}
	return s
}

func (s *Logger) Sync() error { return s.logger.Sync() }
func (s *Logger) GetCore() zapcore.Core {
	return s.logger.Core()
}

// 確保 Logger 實現 logger.Logger 介面
var _ logger.Logger = (*Logger)(nil)

// 確保 Logger 實現 logger.Core 介面
var _ logger.Core = (*Logger)(nil)
