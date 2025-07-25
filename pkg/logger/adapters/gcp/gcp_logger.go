// Package gcp 提供基於 Zap 的 Google Cloud Platform 日誌服務實現
// 將日誌輸出到 Google Cloud Logging 平台
package gcp

import (
	"context"
	"github.com/tomoffice/go-clean-architecture/pkg/logger"

	"cloud.google.com/go/logging"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 定義 GCP Logger 的配置參數
type Config struct {
	ProjectID string       // GCP 專案 ID
	LogName   string       // Cloud Logging 的日誌名稱
	Level     logger.Level // 日誌最低輸出等級
}

// NewDefaultConfig 創建預設配置的 GCP Logger
// 預設配置：指定專案 ID + app-logs 日誌名稱 + Info 等級
func NewDefaultConfig(projectID string) Config {
	return Config{
		ProjectID: projectID,
		LogName:   "app-logs",
		Level:     logger.InfoLevel, // 預設 Info 等級
	}
}

// Logger 實現 logger.Logger 介面，將日誌輸出到 Google Cloud Logging
type Logger struct {
	core   zapcore.Core
	logger *zap.Logger
}

// NewLogger 創建新的 GCP Logger 實例
func NewLogger(cfg Config) (*Logger, error) {
	ctx := context.Background()

	// 1. 建立 Cloud Logging 客戶端
	client, err := logging.NewClient(ctx, cfg.ProjectID)
	if err != nil {
		return nil, err
	}

	// 2. 設定 GCP 日誌寫入器
	gcpWriter := zapcore.AddSync(client.Logger(cfg.LogName).StandardLogger(mapSeverity(cfg.Level)).Writer())

	// 3. 建立統一的編碼器配置（與 Console Logger 保持一致）
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
	enc := zapcore.NewJSONEncoder(encCfg)

	// 4. 建立 Core 和 Logger 實例
	core := zapcore.NewCore(enc, gcpWriter, cfg.Level)
	zl := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	return &Logger{
		core:   core,
		logger: zl,
	}, nil
}
func NewDefaultLogger(projectID string) (logger.Logger, error) {
	return NewLogger(NewDefaultConfig(projectID))
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
	return &Logger{
		core:   l.core,
		logger: l.logger.With(fields...),
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

// GetCore 回傳底層 zapcore.Core，供 MultiLogger 組合使用
func (l *Logger) GetCore() zapcore.Core {
	return l.core
}

// 確保 Logger 實現 logger.Logger 介面
var _ logger.Logger = (*Logger)(nil)

// 確保 Logger 實現 logger.Core 介面
var _ logger.Core = (*Logger)(nil)
