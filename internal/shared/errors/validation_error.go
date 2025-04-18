// Package errors 輸入驗證錯誤（form/JSON 驗證等）
package errors

import "errors"

var (
	ErrInvalidJSON        = errors.New("invalid json body")
	ErrValidationFailed   = errors.New("validation failed")
	ErrMissingField       = errors.New("missing required field")
	ErrInvalidQueryParams = errors.New("invalid query parameters")
)
