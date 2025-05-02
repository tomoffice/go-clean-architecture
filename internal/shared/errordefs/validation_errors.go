// Package errordefs 輸入驗證錯誤（form/JSON 驗證等）
package errordefs

import "errors"

var (
	ErrValidationFailed = errors.New("field validation failed")
	ErrMissingField     = errors.New("missing required field")
)
