package usecase

import (
	"errors"
	infradberrors "module-clean/internal/member/infrastructure/persistence/sqlx"
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
	switch {
	case errors.Is(err, infradberrors.ErrDBRecordNotFound):
		return ErrMemberNotFound
	case errors.Is(err, infradberrors.ErrDBDuplicateKey):
		return ErrMemberAlreadyExists
	case errors.Is(err, infradberrors.ErrDBUpdateNoEffect):
		return ErrMemberUpdateFailed
	case errors.Is(err, infradberrors.ErrDBDeleteNoEffect):
		return ErrMemberDeleteFailed
		// 多種錯誤皆 map 成一種意義
	case errors.Is(err, infradberrors.ErrDBContextTimeout),
		errors.Is(err, infradberrors.ErrDBContextCanceled),
		errors.Is(err, infradberrors.ErrDBConnectionClosed),
		errors.Is(err, infradberrors.ErrDBTransactionDone):
		return ErrMemberDBFailure
	default:
		return ErrUnexpectedMemberUseCaseFail
	}
}
