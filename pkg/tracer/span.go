package tracer

// Span 定義追蹤 span 的基礎介面

//go:generate mockgen -source=span.go -destination=mock/mock_span.go -package=mock
type Span interface {
	End()
}
