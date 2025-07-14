package http

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/tomoffice/go-clean-architecture/internal/shared/errorcode"
	sharederrors "github.com/tomoffice/go-clean-architecture/internal/shared/errordefs"
)

// MapValidationError 將 validator 驗證失敗錯誤，轉換為 error code 與人類可讀訊息
func MapMemberValidationError(err error) (int, string) {
	var valErr validator.ValidationErrors

	if errors.As(err, &valErr) {
		fieldErr := valErr[0] // 取第一個欄位錯誤回報
		return errorcode.ErrValidationFailed,
			//欄位驗證失敗
			fmt.Sprintf("Column '%s' validation failed (Rule: %s)", fieldErr.Field(), fieldErr.ActualTag())
	}

	// fallback，理論上不應該到這裡
	return errorcode.ErrValidationFailed, sharederrors.ErrValidationFailed.Error()
}
