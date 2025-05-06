package usecase

import (
	"errors"
	sqlite2 "module-clean/internal/modules/member/infrastructure/persistence/sqlx/sqlite"
)

// MapInfraErrorToUseCaseError 將 infra 錯誤轉換為 usecase 層語意錯誤
func MapInfraErrorToUseCaseError(err error) error {
	if err == nil {
		return nil
	}

	// 優先比對 CustomError
	switch {
	case errors.Is(err, sqlite2.ErrDBRecordNotFound):
		return ErrMemberNotFound
	case errors.Is(err, sqlite2.ErrDBDuplicateKey):
		return ErrMemberAlreadyExists
	case errors.Is(err, sqlite2.ErrDBUpdateNoEffect):
		return ErrMemberUpdateFailed
	case errors.Is(err, sqlite2.ErrDBDeleteNoEffect):
		return ErrMemberDeleteFailed
	}

	// 萬一錯誤不是 CustomError instance，要用 As 抓 DBError
	var dbErr *sqlite2.DBError
	if errors.As(err, &dbErr) {
		switch dbErr.CustomError {
		case sqlite2.ErrDBContextTimeout,
			sqlite2.ErrDBContextCanceled,
			sqlite2.ErrDBConnectionClosed,
			sqlite2.ErrDBTransactionDone,
			sqlite2.ErrDBUnknown:
			return ErrMemberDBFailure
		}
	}

	// fallback 落到底了，真的不認得就回 Unexpected
	return ErrUnexpectedMemberUseCaseFail
}
