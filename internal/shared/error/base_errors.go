// Package error 系統通用錯誤（系統、不可預期）
package error

import "errors"

var (
	ErrInternalServer = errors.New("internal server error")
	ErrUnavailable    = errors.New("service temporarily unavailable")
)
