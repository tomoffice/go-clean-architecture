//go:generate mockgen -source=tracer.go -destination=mock/mock_tracer.go -package=mock

// Package tracer 提供分散式追蹤功能的封裝
// 基於 OpenTelemetry 實作，隱藏底層複雜性，提供清晰的 API
package tracer

import (
	"context"
)

// Tracer 定義追蹤器的核心介面
type Tracer interface {
	Start(ctx context.Context, name string) (context.Context, Span)
}
