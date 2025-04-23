package http

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"module-clean/internal/member/usecase"
	sharederrors "module-clean/internal/shared/error"
	"module-clean/internal/shared/errorcode"
)

func MapErrorToResponse(err error) (int, string) {
	if code, msg := mapValidationError(err); code != 0 {
		return code, msg
	}
	if code, msg := mapMemberUseCaseError(err); code != 0 {
		return code, msg
	}
	return errorcode.ErrInternalServer, sharederrors.ErrInternalServer.Error()
}
func mapValidationError(err error) (int, string) {
	var typeErr *json.UnmarshalTypeError
	var syntaxErr *json.SyntaxError
	var valErr validator.ValidationErrors

	switch {
	case errors.As(err, &typeErr):
		return errorcode.ErrInvalidJSONType, "欄位類型錯誤：" + typeErr.Field
	case errors.As(err, &syntaxErr):
		return errorcode.ErrInvalidJSONSyntax, "JSON 格式錯誤：" + syntaxErr.Error()
	case errors.As(err, &valErr):
		return errorcode.ErrValidationFailed, "欄位驗證失敗：" + valErr.Error()
	default:
		return errorcode.ErrInvalidJSONInput, sharederrors.ErrInvalidJSONInput.Error()
	}
}
func mapMemberUseCaseError(err error) (int, string) {
	switch {
	case errors.Is(err, usecase.ErrMemberNotFound):
		return errorcode.ErrMemberNotFound, usecase.ErrMemberNotFound.Error()
	case errors.Is(err, usecase.ErrMemberAlreadyExists):
		return errorcode.ErrMemberAlreadyExists, usecase.ErrMemberAlreadyExists.Error()
	case errors.Is(err, usecase.ErrMemberUpdateFailed):
		return errorcode.ErrMemberUpdateFailed, usecase.ErrMemberUpdateFailed.Error()
	case errors.Is(err, usecase.ErrMemberDeleteFailed):
		return errorcode.ErrMemberDeleteFailed, usecase.ErrMemberDeleteFailed.Error()
	case errors.Is(err, usecase.ErrUnexpectedMemberUseCaseFail):
		return errorcode.ErrInternalServer, usecase.ErrUnexpectedMemberUseCaseFail.Error()
	default:
		return 0, ""
	}
}
