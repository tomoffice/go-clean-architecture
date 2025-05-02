package response

import (
	"github.com/gin-gonic/gin"
	"module-clean/internal/shared/enum"
	"strconv"
)

type APIResponse[T any] struct {
	Data  T             `json:"data,omitempty"`
	Error *ErrorPayload `json:"errordefs,omitempty"`
	Meta  *MetaPayload  `json:"meta,omitempty"`
	BaseResponse
}
type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type MetaPayload struct {
	Total  int `json:"total"`
	Page   int `json:"page"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func SuccessAPIResponse[T any](ctx *gin.Context, data T) {
	ctx.JSON(200, APIResponse[T]{
		Data:         data,
		BaseResponse: NewBaseResponse(enum.APIStatusSuccess),
	})
}
func SuccessAPIResponseWithMeta[T any](ctx *gin.Context, data T, total, page, limit, offset int) {
	ctx.JSON(200, APIResponse[T]{
		Data:         data,
		Meta:         &MetaPayload{Total: total, Page: page, Limit: limit, Offset: offset},
		BaseResponse: NewBaseResponse(enum.APIStatusSuccess),
	})
}
func FailureAPIResponse(ctx *gin.Context, statusCode int, errorCode int, message string) {
	ctx.JSON(statusCode, APIResponse[any]{
		Error: &ErrorPayload{
			Code:    strconv.Itoa(errorCode),
			Message: message,
		},
		BaseResponse: NewBaseResponse(enum.APIStatusFailed),
	})
}
