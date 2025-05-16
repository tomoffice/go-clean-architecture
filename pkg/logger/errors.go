package logger

import "errors"

// Logger 專用錯誤（如 ErrNoAdapters）

var (
	ErrNoAdapters = errors.New("logger: no adapters specified")
)
