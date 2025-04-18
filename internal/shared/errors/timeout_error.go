// Package errors 超時、逾期等
package errors

import "errors"

var (
	ErrRequestTimeout = errors.New("request timeout")
	ErrContextTimeout = errors.New("context deadline exceeded")
)
