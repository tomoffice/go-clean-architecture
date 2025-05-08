package response

import (
	"github.com/gin-gonic/gin"
	"module-clean/internal/shared/common/enum"
	"module-clean/internal/shared/interface_adapter/viewmodel/http"
	"strconv"
)

func WriteSuccess[T any](ctx *gin.Context, data T) {
	ctx.JSON(200, http.HTTPResponse[T]{
		Data:             data,
		BaseHTTPResponse: http.NewBaseHTTPResponse(enum.APIStatusSuccess),
	})
}
func WriteSuccessWithMeta[T any](ctx *gin.Context, data T, total, page, limit, offset int) {
	ctx.JSON(200, http.HTTPResponse[T]{
		Data:             data,
		Meta:             &http.MetaPayload{Total: total, Page: page, Limit: limit, Offset: offset},
		BaseHTTPResponse: http.NewBaseHTTPResponse(enum.APIStatusSuccess),
	})
}
func WriteFailure(ctx *gin.Context, statusCode int, errorCode int, message string) {
	ctx.JSON(statusCode, http.HTTPResponse[any]{
		Error: &http.ErrorPayload{
			Code:    strconv.Itoa(errorCode),
			Message: message,
		},
		BaseHTTPResponse: http.NewBaseHTTPResponse(enum.APIStatusFailed),
	})
}
