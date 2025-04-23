// Package error 輸入驗證錯誤（form/JSON 驗證等）
package error

import "errors"

var (
	ErrInvalidJSONSyntax  = errors.New("invalid JSON syntax")
	ErrInvalidJSONType    = errors.New("invalid JSON field type")
	ErrValidationFailed   = errors.New("field validation failed")
	ErrMissingField       = errors.New("missing required field")
	ErrInvalidQueryParams = errors.New("invalid query parameters")
	ErrInvalidJSONInput   = errors.New("invalid JSON input")
)
