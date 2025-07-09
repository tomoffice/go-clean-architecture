// Package seq 提供基於 Zap 和 Logrus 的 Seq 日誌服務實現
// 透過 Bridge 模式將 Zap 日誌轉發到 Seq 日誌平台
package seq

import (
	"time"

	"github.com/tomoffice/module-clean/pkg/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Config 定義 Seq Logger 的配置參數
type Config struct {
	Endpoint string        // Seq 服務端點，例如：http://seq:5341
	APIKey   string        // Seq API 金鑰，可選參數
	Level    zapcore.Level // 日誌最低輸出等級
}

// NewDefaultConfig 創建預設配置的 Seq Logger
// 預設配置：本地端點 localhost:5341 + Info 等級
func NewDefaultConfig() (logger.Logger, error) {
	return NewLogger(Config{
		Endpoint: "http://localhost:5341",
		APIKey:   "",
		Level:    zapcore.InfoLevel,
	})
}

// Logger 實現 logger.Logger 介面，提供 Seq 日誌輸出功能
type Logger struct {
	core   zapcore.Core
	logger *zap.Logger
}

// NewLogger 創建新的 Seq Logger 實例
func NewLogger(cfg Config) (logger.Logger, error) {
	sender := NewLogrusSender(cfg)

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
	return &Logger{core: core, logger: zl}, nil
}

// Debug 輸出 Debug 等級的日誌訊息
func (s *Logger) Debug(msg string, field ...logger.Field) { s.logger.Debug(msg, field...) }

// Info 輸出 Info 等級的日誌訊息
func (s *Logger) Info(msg string, field ...logger.Field) { s.logger.Info(msg, field...) }

// Warn 輸出 Warn 等級的日誌訊息
func (s *Logger) Warn(msg string, field ...logger.Field) { s.logger.Warn(msg, field...) }

// Error 輸出 Error 等級的日誌訊息
func (s *Logger) Error(msg string, field ...logger.Field) { s.logger.Error(msg, field...) }

// With 返回一個新的 Logger 實例，預先附加指定的結構化欄位
func (s *Logger) With(field ...logger.Field) logger.Logger {
	return &Logger{
		core:   s.core.With(field),
		logger: s.logger.With(field...),
	}
}

func (s *Logger) Sync() error { return s.logger.Sync() }
func (s *Logger) GetCore() zapcore.Core {
	return s.logger.Core()
}
