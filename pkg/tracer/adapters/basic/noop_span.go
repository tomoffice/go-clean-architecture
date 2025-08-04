package basic

import (
	"github.com/tomoffice/go-clean-architecture/pkg/tracer"
)

// NoOpSpan 是一個不執行任何操作的 span 實作
// 用於當 tracer 被停用時
type NoOpSpan struct{}

// End 實作 tracer.Span.End 方法，但不執行任何操作
func (s *NoOpSpan) End() {
	// 不執行任何操作
}

// 確保實作正確的介面
var _ tracer.Span = (*NoOpSpan)(nil)
