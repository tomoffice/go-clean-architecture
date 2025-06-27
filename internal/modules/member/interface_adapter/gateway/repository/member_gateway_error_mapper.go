package repository

import (
	"errors"
	"fmt"
	"module-clean/internal/modules/member/driver/persistence/sqlx/sqlite"
	"module-clean/internal/modules/member/usecase"
)

// MapInfraErrorToUsecaseError 將底層 infra 錯誤直接轉為 usecase 定義的 sentinel error
func MapInfraErrorToUsecaseError(err error) error {
	if err == nil {
		return nil
	}
	// 先比對 CustomError
	switch {
	case errors.Is(err, sqlite.ErrDBRecordNotFound):
		return usecase.ErrUseCaseMemberNotFound
	case errors.Is(err, sqlite.ErrDBDuplicateKey):
		return usecase.ErrUseCaseMemberAlreadyExists
	case errors.Is(err, sqlite.ErrDBNoEffect):
		return usecase.ErrUseCaseMemberNoEffect
	}
	// 再處理 DBError 類型
	var dbErr *sqlite.DBError
	if errors.As(err, &dbErr) {
		return fmt.Errorf("%w: %v", usecase.ErrUseCaseMemberDBError, dbErr.RawError)
	}
	// fallback：其他未知錯誤
	return usecase.ErrUseCaseMemberUnexpectedError
}
