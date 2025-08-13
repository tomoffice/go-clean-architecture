package basic

import (
	"context"

	"github.com/tomoffice/go-clean-architecture/pkg/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Config 定義 Basic tracer 的配置
type Config struct {
	// ServiceName 服務名稱，用於識別不同的服務
	ServiceName string
	// Enabled 是否啟用追蹤功能
	Enabled bool
}

// DefaultConfig 返回預設配置
func DefaultConfig() Config {
	return Config{
		ServiceName: "unknown-service",
		Enabled:     true,
	}
}

// NewConfig 創建新的配置
func NewConfig(serviceName string, enabled bool) Config {
	return Config{
		ServiceName: serviceName,
		Enabled:     enabled,
	}
}

// Tracer 實作 tracer.Tracer 介面，基於 OpenTelemetry 介面但不需要完整初始化
type Tracer struct {
	tracer      trace.Tracer
	serviceName string
	enabled     bool
}

// NewTracer 創建基於 OpenTelemetry 介面的基本 tracer 實例
func NewTracer(cfg Config) tracer.Tracer {
	return &Tracer{
		tracer:      otel.Tracer(cfg.ServiceName),
		serviceName: cfg.ServiceName,
		enabled:     cfg.Enabled,
	}
}

// Start 實作 tracer.Tracer.Start 方法
func (t *Tracer) Start(ctx context.Context, name string) (context.Context, tracer.Span) {
	if !t.enabled {
		// 如果未啟用，返回 no-op span
		return ctx, &NoOpSpan{}
	}
	
	ctx, span := t.tracer.Start(ctx, name)
	return ctx, NewSpan(span)
}


// 確保 Tracer 實作 tracer.Tracer 介面
var _ tracer.Tracer = (*Tracer)(nil)