package adapter

import (
	"github.com/gin-gonic/gin"
	httpx "github.com/tomoffice/go-clean-architecture/internal/interface_adapter/transport/http"
)

type ginRouter struct{ grp *gin.RouterGroup }

func NewRouter(grp *gin.RouterGroup) httpx.Router { return ginRouter{grp: grp} }

func (r ginRouter) GET(p string, h httpx.HandlerFunc)    { r.grp.GET(p, wrap(h)) }
func (r ginRouter) POST(p string, h httpx.HandlerFunc)   { r.grp.POST(p, wrap(h)) }
func (r ginRouter) PUT(p string, h httpx.HandlerFunc)    { r.grp.PUT(p, wrap(h)) }
func (r ginRouter) PATCH(p string, h httpx.HandlerFunc)  { r.grp.PATCH(p, wrap(h)) }
func (r ginRouter) DELETE(p string, h httpx.HandlerFunc) { r.grp.DELETE(p, wrap(h)) }
