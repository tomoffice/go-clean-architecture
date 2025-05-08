package errordefs

import (
	"encoding/json"
	"errors"
	"fmt"
	"module-clean/internal/shared/common/errorcode"
)

// MapGinBindingError 將 binding 時產生的錯誤轉換為 (error code, message)
func MapGinBindingError(err error) (int, string) {
	var (
		syntaxErr *json.SyntaxError
		typeErr   *json.UnmarshalTypeError
	)

	switch {
	case errors.As(err, &syntaxErr):
		return errorcode.ErrInvalidJSONSyntax,
			fmt.Sprintf("%s: %s", ErrInvalidJSONSyntax.Error(), syntaxErr.Error())

	case errors.As(err, &typeErr):
		return errorcode.ErrInvalidJSONType,
			fmt.Sprintf("%s: field=%s", ErrInvalidJSONType.Error(), typeErr.Field)
	}

	switch {
	case errors.Is(err, ErrFromJSON):
		return errorcode.ErrInvalidJSONInput, ErrFromJSON.Error()
	case errors.Is(err, ErrFromQuery):
		return errorcode.ErrInvalidQueryParams, ErrFromQuery.Error()
	case errors.Is(err, ErrFromURI):
		return errorcode.ErrInvalidURIParams, ErrFromURI.Error()
	case errors.Is(err, ErrFromForm):
		return errorcode.ErrInvalidFormData, ErrFromForm.Error()
	case errors.Is(err, ErrFromHeader):
		return errorcode.ErrInvalidHeaderParams, ErrFromHeader.Error()

	}

	// fallback
	return errorcode.ErrInternalServer, ErrUnexpectedBinding.Error()
}
