package errordefs

import "errors"

// 輸入來源標記錯誤
var (
	ErrFromJSON   = errors.New("input from JSON body")
	ErrFromQuery  = errors.New("input from query string")
	ErrFromForm   = errors.New("input from form data")
	ErrFromURI    = errors.New("input from URI path")
	ErrFromHeader = errors.New("input from header")
)

// Bind 發生錯誤的型別
var (
	ErrInvalidJSONSyntax   = errors.New("invalid JSON syntax")
	ErrInvalidJSONType     = errors.New("invalid JSON field type")
	ErrUnexpectedBinding   = errors.New("unexpected binding error")
	ErrInvalidQueryParams  = errors.New("invalid query parameters")
	ErrInvalidFormData     = errors.New("invalid form data")
	ErrInvalidURIParams    = errors.New("invalid URI parameters")
	ErrInvalidHeaderParams = errors.New("invalid header parameters")
)
