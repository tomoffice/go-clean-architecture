// Package errordefs 超時、逾期等
package errordefs

import "errors"

var (
	ErrRequestTimeout  = errors.New("request timeout")
	ErrContextTimeout  = errors.New("context deadline exceeded")
	ErrContextCanceled = errors.New("context canceled")
)
