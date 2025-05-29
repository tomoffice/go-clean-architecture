package usecase

import (
	"errors"
)

// MemberUseCase 錯誤碼
var (
	ErrUseCaseMemberNotFound       = errors.New("usecase: member not found")
	ErrUseCaseMemberAlreadyExists  = errors.New("usecase: email already exists")
	ErrUseCaseMemberUpdateFailed   = errors.New("usecase: member update failed")
	ErrUseCaseMemberDeleteFailed   = errors.New("usecase: member delete failed")
	ErrUseCaseMemberDBFailure      = errors.New("usecase: member db operation failed")
	ErrUseCaseMemberUnexpectedFail = errors.New("usecase: unexpected member usecase errordefs")
	ErrUseCaseGatewayFailure       = errors.New("usecase: member gateway failure")
)
