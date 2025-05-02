// Package errordefs 系統通用錯誤（系統、不可預期）
package errordefs

import "errors"

var (
	ErrInternalServer = errors.New("internal server errordefs")
	ErrUnavailable    = errors.New("service temporarily unavailable")
)
