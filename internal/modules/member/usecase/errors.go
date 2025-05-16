package usecase

import (
	"errors"
)

// MemberUseCase 錯誤碼
var (
	ErrMemberNotFound              = errors.New("usecase: member not found")
	ErrMemberAlreadyExists         = errors.New("usecase: email already exists")
	ErrMemberUpdateFailed          = errors.New("usecase: member update failed")
	ErrMemberDeleteFailed          = errors.New("usecase: member delete failed")
	ErrMemberDBFailure             = errors.New("usecase: member db operation failed")
	ErrUnexpectedMemberUseCaseFail = errors.New("usecase: unexpected member usecase errordefs")
	ErrGatewayFailure              = errors.New("usecase: member gateway failure")
)
