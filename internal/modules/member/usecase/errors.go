package usecase

import (
	"errors"
)

// MemberUseCase 錯誤碼
var (
	ErrUseCaseMemberNotFound      = errors.New("usecase: member not found")
	ErrUseCaseMemberAlreadyExists = errors.New("usecase: member already exists")
	ErrUseCaseMemberUpdateFailed  = errors.New("usecase: member update failed")
	ErrUseCaseMemberDeleteFailed  = errors.New("usecase: member delete failed")
	ErrUseCaseMemberDBError       = errors.New("usecase: member db operation error")
	ErrUseCaseUnexpectedError     = errors.New("usecase: member usecase unexpected error")
	ErrUseCaseGatewayError        = errors.New("usecase: member gateway error")
)
