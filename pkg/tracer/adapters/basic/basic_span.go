package basic

import (
	"github.com/tomoffice/go-clean-architecture/pkg/tracer"
	"go.opentelemetry.io/otel/trace"
)

// Span 實作 tracer.Span 介面，基於 OpenTelemetry
type Span struct {
	span trace.Span // 直接包裝 OpenTelemetry 的原始 span
}

// NewSpan 創建新的 Span 實例
func NewSpan(span trace.Span) tracer.Span {
	return &Span{
		span: span,
	}
}

// End 實作 tracer.Span.End 方法
func (s *Span) End() {
	s.span.End()
}

// NoOpSpan 是一個不執行任何操作的 span 實作
// 用於當 tracer 被停用時
type NoOpSpan struct{}

// End 實作 tracer.Span.End 方法，但不執行任何操作
func (s *NoOpSpan) End() {
	// 不執行任何操作
}

// 確保實作正確的介面
var _ tracer.Span = (*Span)(nil)
var _ tracer.Span = (*NoOpSpan)(nil)