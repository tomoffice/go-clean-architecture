package tracer

// Span 定義追蹤 span 的基礎介面
type Span interface {
	End()
}
