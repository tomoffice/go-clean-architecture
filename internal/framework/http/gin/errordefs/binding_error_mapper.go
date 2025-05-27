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
			//JSON 格式錯誤
			fmt.Sprintf("%s: %s", ErrInvalidJSONSyntax.Error(), syntaxErr.Error())

	case errors.As(err, &typeErr):
		return errorcode.ErrInvalidJSONType,
			//JSON 欄位類型錯誤
			fmt.Sprintf("%s: field=%s", ErrInvalidJSONType.Error(), typeErr.Field)
	}

	// fallback: 其它一律回 InternalServer 與通用訊息 "參數格式錯誤"
	return errorcode.ErrInvalidParams, ErrInvalidParams.Error()
}
