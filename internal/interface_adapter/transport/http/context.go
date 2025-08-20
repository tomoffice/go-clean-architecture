package http

import (
	"context"
	"net/http"
)

//go:generate mockgen -source=context.go -destination=../../controller/mock/mock_member_http_context.go -package=mock
type Context interface {
	// 綁定
	BindJSON(dest any) error
	BindQuery(dest any) error
	BindURI(dest any) error

	// 讀取
	GetHeader(key string) string
	RequestCtx() context.Context

	// 回傳原生*http.Request
	Request() *http.Request // 如果需要原生的 *http.Request，可以加上這個方法

	// 回應
	Header(key, val string)
	Status(code int)
	JSON(code int, body any)
}
