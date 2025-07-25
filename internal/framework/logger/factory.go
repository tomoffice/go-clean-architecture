package logger

import (
	"context"
	"fmt"

	"github.com/tomoffice/go-clean-architecture/config"
	"github.com/tomoffice/go-clean-architecture/pkg/logger"
	"github.com/tomoffice/go-clean-architecture/pkg/logger/adapters/console"
	"github.com/tomoffice/go-clean-architecture/pkg/logger/adapters/gcp"
	"github.com/tomoffice/go-clean-architecture/pkg/logger/adapters/seq"
	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// NewLogger 根據配置創建 logger 實例並初始化 OpenTelemetry
// 返回組合後的 logger 和 cleanup 函數
func NewLogger(cfg config.LoggerConfig) (logger.Logger, func() error, error) {
	// 1. 初始化 OpenTelemetry（目前預設啟用）
	cleanup := func() error { return nil }

	// 創建 TracerProvider
	tp := sdktrace.NewTracerProvider()

	// 設定全域 TracerProvider
	otel.SetTracerProvider(tp)

	cleanup = func() error {
		return tp.Shutdown(context.Background())
	}

	// 2. 創建 logger adapters
	var loggers []logger.Logger

	// Console Logger
	if cfg.Console.Enabled {
		consoleLogger, err := console.NewLogger(console.Config{
			Level:  logger.ParseLevel(cfg.Console.Level),
			Format: logger.ParseFormat(cfg.Console.Format),
		})
		if err != nil {
			return nil, cleanup, fmt.Errorf("創建 Console Logger 失敗: %w", err)
		}
		loggers = append(loggers, consoleLogger)
	}

	// GCP Logger
	if cfg.GCP.Enabled {
		gcpLogger, err := gcp.NewLogger(gcp.Config{
			ProjectID: cfg.GCP.ProjectID,
			LogName:   "app-log", // 使用預設值，之後可以加到配置中
			Level:     logger.ParseLevel(cfg.GCP.Level),
		})
		if err != nil {
			return nil, cleanup, fmt.Errorf("創建 GCP Logger 失敗: %w", err)
		}
		loggers = append(loggers, gcpLogger)
	}

	// Seq Logger
	if cfg.Seq.Enabled {
		seqLogger, err := seq.NewLogger(seq.Config{
			Endpoint:             cfg.Seq.Endpoint,
			APIKey:               cfg.Seq.APIKey,
			Level:                logger.ParseLevel(cfg.Seq.Level),
			ConsoleOutputEnabled: cfg.Seq.ConsoleOutputEnabled,
		})
		if err != nil {
			return nil, cleanup, fmt.Errorf("創建 Seq Logger 失敗: %w", err)
		}
		loggers = append(loggers, seqLogger)
	}

	// 3. 如果沒有啟用任何 logger，創建預設的 console logger
	if len(loggers) == 0 {
		defaultLogger, err := console.NewDefaultLogger()
		if err != nil {
			return nil, cleanup, fmt.Errorf("創建預設 Console Logger 失敗: %w", err)
		}
		loggers = append(loggers, defaultLogger)
	}

	// 4. 使用 teeLogger 組合多個 logger
	var appLogger logger.Logger
	if len(loggers) == 1 {
		// 只有一個 logger，直接使用
		appLogger = loggers[0]
	} else {
		// 多個 logger，使用 tee 組合
		appLogger = logger.NewTeeLogger(loggers...)
	}

	return appLogger, cleanup, nil
}
