package repository

import (
	"errors"
	"fmt"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/framework/persistence/sqlx/mcsqlite"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/usecase"
)

// MapInfraErrorToUsecaseError 將底層 infra 錯誤直接轉為 usecase 定義的 sentinel error
func MapInfraErrorToUsecaseError(err error) error {
	if err == nil {
		return nil
	}
	// 先處理 gateway 層的業務語意錯誤
	if errors.Is(err, ErrGatewayMemberMappingError) {
		return usecase.ErrMemberMappingError
	}
	// 先比對 CustomError
	switch {
	case errors.Is(err, mcsqlite.ErrDBRecordNotFound):
		return usecase.ErrMemberNotFound
	case errors.Is(err, mcsqlite.ErrDBDuplicateKey):
		return usecase.ErrMemberAlreadyExists
	case errors.Is(err, mcsqlite.ErrDBNoEffect):
		return usecase.ErrMemberNoEffect
	}
	// 再處理 DBError 類型
	var dbErr *mcsqlite.DBError
	if errors.As(err, &dbErr) {
		return fmt.Errorf("%w: %v", usecase.ErrMemberDBError, dbErr.RawError)
	}
	// fallback：其他未知錯誤
	return usecase.ErrMemberUnexpectedError
}
