package repository

import (
	"errors"
	"fmt"
	"module-clean/internal/modules/member/driver/persistence/sqlx/sqlite"
)

// MapInfraErrorToGatewayError 將 infra 錯誤轉換為 gateway 層錯誤（並保留原始 error）
func MapInfraErrorToGatewayError(err error) error {
	if err == nil {
		return nil
	}
	// 優先比對 CustomError
	switch {
	case errors.Is(err, sqlite.ErrDBRecordNotFound):
		return ErrGatewayMemberNotFound
	case errors.Is(err, sqlite.ErrDBDuplicateKey):
		return ErrGatewayMemberAlreadyExists
	case errors.Is(err, sqlite.ErrDBNoEffect):
		return ErrGatewayMemberNoEffect
	}
	// 萬一錯誤不是 CustomError instance，要用 As 抓 DBError
	var dbErr *sqlite.DBError
	if errors.As(err, &dbErr) {
		if errors.Is(dbErr.CustomError, sqlite.ErrDBContextTimeout) ||
			errors.Is(dbErr.CustomError, sqlite.ErrDBContextCanceled) ||
			errors.Is(dbErr.CustomError, sqlite.ErrDBConnectionClosed) ||
			errors.Is(dbErr.CustomError, sqlite.ErrDBTransactionDone) ||
			errors.Is(dbErr.CustomError, sqlite.ErrDBUnexpectedError) {
			return fmt.Errorf("%w: %w", ErrGatewayMemberDBError, dbErr.RawError)
		}
	}
	// fallback：未知錯誤也包裝
	return fmt.Errorf("gateway: %w: %w", ErrGatewayMemberUnexpectedError, err)
}
