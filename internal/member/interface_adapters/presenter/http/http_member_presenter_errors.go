package http

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/validator/v10"
	"module-clean/internal/member/usecase"
	sharederrors "module-clean/internal/shared/error"
	"module-clean/internal/shared/errorcode"
)

func MapInputValidationError(err error) (int, string) {
	var (
		typeErr   *json.UnmarshalTypeError
		syntaxErr *json.SyntaxError
		valErr    validator.ValidationErrors
	)

	isFromJSON := errors.Is(err, sharederrors.ErrFromJSON)
	isFromQuery := errors.Is(err, sharederrors.ErrFromQuery)
	isFromURI := errors.Is(err, sharederrors.ErrFromURI)
	isFromForm := errors.Is(err, sharederrors.ErrFromForm)
	isFromHeader := errors.Is(err, sharederrors.ErrFromHeader)

	switch {
	case errors.As(err, &typeErr):
		return errorcode.ErrInvalidJSONType, "JSON 欄位型別錯誤：" + typeErr.Field
	case errors.As(err, &syntaxErr):
		return errorcode.ErrInvalidJSONSyntax, "JSON 格式錯誤：" + syntaxErr.Error()
	case errors.As(err, &valErr):
		switch {
		case isFromQuery:
			return errorcode.ErrInvalidQueryParams, "查詢參數驗證失敗：" + valErr.Error()
		case isFromURI:
			return errorcode.ErrInvalidURIParams, "URI 參數驗證失敗：" + valErr.Error()
		case isFromForm:
			return errorcode.ErrInvalidFormData, "表單資料驗證失敗：" + valErr.Error()
		case isFromHeader:
			return errorcode.ErrInvalidHeaderParams, "Header 參數驗證失敗：" + valErr.Error()
		case isFromJSON:
			return errorcode.ErrValidationFailed, "JSON 欄位驗證失敗：" + valErr.Error()
		default:
			return errorcode.ErrValidationFailed, "輸入欄位驗證失敗：" + valErr.Error()
		}
	default:
		switch {
		case isFromQuery:
			return errorcode.ErrInvalidQueryParams, "查詢參數格式錯誤"
		case isFromURI:
			return errorcode.ErrInvalidURIParams, "URI 參數格式錯誤"
		case isFromForm:
			return errorcode.ErrInvalidFormData, "表單資料格式錯誤"
		case isFromHeader:
			return errorcode.ErrInvalidHeaderParams, "Header 格式錯誤"
		case isFromJSON:
			return errorcode.ErrInvalidJSONInput, "JSON 格式錯誤"
		default:
			return errorcode.ErrInternalServer, "未知輸入錯誤"
		}
	}
}
func MapMemberUseCaseError(err error) (int, string) {
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
		return errorcode.ErrUnexpectedMemberUseCaseFail, usecase.ErrUnexpectedMemberUseCaseFail.Error()
	case errors.Is(err, usecase.ErrMemberDBFailure):
		return errorcode.ErrMemberDBFailure, usecase.ErrMemberDBFailure.Error()
	default:
		return errorcode.ErrInternalServer, sharederrors.ErrInternalServer.Error()
	}
}
