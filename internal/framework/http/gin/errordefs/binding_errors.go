package errordefs

import "errors"

// 主要定義錯誤型別，不需要標記來源
var (
	ErrInvalidJSONSyntax = errors.New("invalid JSON syntax")
	ErrInvalidJSONType   = errors.New("invalid JSON field type")
	ErrInvalidParams     = errors.New("invalid parameters") // 通用參數錯誤，包含 Query/Form/URI/Header 等
)
