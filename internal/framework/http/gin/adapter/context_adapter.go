package adapter

import (
	"context"
	memberhttp "github.com/tomoffice/go-clean-architecture/internal/interface_adapter/transport/http"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ginContext struct{ c *gin.Context }

// 綁定（沿用 Gin 的 ShouldBindXXX）
func (g ginContext) BindJSON(v any) error  { return g.c.ShouldBindJSON(v) }
func (g ginContext) BindQuery(v any) error { return g.c.ShouldBindQuery(v) }
func (g ginContext) BindURI(v any) error   { return g.c.ShouldBindUri(v) }

// 讀取
func (g ginContext) GetHeader(k string) string   { return g.c.GetHeader(k) }
func (g ginContext) RequestCtx() context.Context { return g.c.Request.Context() } // ★ 關鍵
func (g ginContext) Request() *http.Request      { return g.c.Request }

// 回應
func (g ginContext) Header(k, v string)      { g.c.Header(k, v) }
func (g ginContext) Status(code int)         { g.c.Status(code) }
func (g ginContext) JSON(code int, body any) { g.c.JSON(code, body) }

// 包裝 handler
func wrap(h memberhttp.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) { h(ginContext{c}) }
}
