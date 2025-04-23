package usecase

import (
	"errors"
	infradberrors "module-clean/internal/member/infrastructure/persistence/sqlx"
)

var (
	ErrMemberNotFound              = errors.New("member not found")
	ErrMemberAlreadyExists         = errors.New("email already exists")
	ErrMemberUpdateFailed          = errors.New("member update failed")
	ErrMemberDeleteFailed          = errors.New("member delete failed")
	ErrUnexpectedMemberUseCaseFail = errors.New("unexpected member usecase error")
)

// MapInfraErrorToUseCaseError 將 infra 錯誤轉換為 usecase 層語意錯誤
func MapInfraErrorToUseCaseError(err error) error {
	switch {
	case errors.Is(err, infradberrors.ErrDBRecordNotFound):
		return ErrMemberNotFound
	case errors.Is(err, infradberrors.ErrDBDuplicateKey):
		return ErrMemberAlreadyExists
	case errors.Is(err, infradberrors.ErrDBUpdateNoEffect):
		return ErrMemberUpdateFailed
	case errors.Is(err, infradberrors.ErrDBDeleteNoEffect):
		return ErrMemberDeleteFailed
	default:
		return ErrUnexpectedMemberUseCaseFail
	}
}
