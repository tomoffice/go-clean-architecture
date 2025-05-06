package response

import (
	"github.com/gin-gonic/gin"
	"module-clean/internal/shared/enum"
	sharedresponse "module-clean/internal/shared/response"
	"strconv"
)

func SuccessAPIResponse[T any](ctx *gin.Context, data T) {
	ctx.JSON(200, sharedresponse.APIResponse[T]{
		Data:         data,
		BaseResponse: sharedresponse.NewBaseResponse(enum.APIStatusSuccess),
	})
}
func SuccessAPIResponseWithMeta[T any](ctx *gin.Context, data T, total, page, limit, offset int) {
	ctx.JSON(200, sharedresponse.APIResponse[T]{
		Data:         data,
		Meta:         &sharedresponse.MetaPayload{Total: total, Page: page, Limit: limit, Offset: offset},
		BaseResponse: sharedresponse.NewBaseResponse(enum.APIStatusSuccess),
	})
}
func FailureAPIResponse(ctx *gin.Context, statusCode int, errorCode int, message string) {
	ctx.JSON(statusCode, sharedresponse.APIResponse[any]{
		Error: &sharedresponse.ErrorPayload{
			Code:    strconv.Itoa(errorCode),
			Message: message,
		},
		BaseResponse: sharedresponse.NewBaseResponse(enum.APIStatusFailed),
	})
}
