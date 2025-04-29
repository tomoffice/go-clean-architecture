package usecase

import (
	"errors"
	"module-clean/internal/member/infrastructure/persistence/sqlx/sqlite"
)

// MemberUseCase 錯誤碼
var (
	ErrMemberNotFound              = errors.New("usecase: member not found")
	ErrMemberAlreadyExists         = errors.New("usecase: email already exists")
	ErrMemberUpdateFailed          = errors.New("usecase: member update failed")
	ErrMemberDeleteFailed          = errors.New("usecase: member delete failed")
	ErrMemberDBFailure             = errors.New("usecase: member db operation failed")
	ErrUnexpectedMemberUseCaseFail = errors.New("usecase: unexpected member usecase error")
)

// MapInfraErrorToUseCaseError 將 infra 錯誤轉換為 usecase 層語意錯誤
func MapInfraErrorToUseCaseError(err error) error {
	if err == nil {
		return nil
	}

	// 優先比對 CustomError
	switch {
	case errors.Is(err, sqlite.ErrDBRecordNotFound):
		return ErrMemberNotFound
	case errors.Is(err, sqlite.ErrDBDuplicateKey):
		return ErrMemberAlreadyExists
	case errors.Is(err, sqlite.ErrDBUpdateNoEffect):
		return ErrMemberUpdateFailed
	case errors.Is(err, sqlite.ErrDBDeleteNoEffect):
		return ErrMemberDeleteFailed
	}

	// 萬一錯誤不是 CustomError instance，要用 As 抓 DBError
	var dbErr *sqlite.DBError
	if errors.As(err, &dbErr) {
		switch dbErr.CustomError {
		case sqlite.ErrDBContextTimeout,
			sqlite.ErrDBContextCanceled,
			sqlite.ErrDBConnectionClosed,
			sqlite.ErrDBTransactionDone,
			sqlite.ErrDBUnknown:
			return ErrMemberDBFailure
		}
	}

	// fallback 落到底了，真的不認得就回 Unexpected
	return ErrUnexpectedMemberUseCaseFail
}
