package usecase

import (
	"errors"
	gateway "module-clean/internal/modules/member/interface_adapter/gateway/repository"
)

// MapInfraErrorToUseCaseError 將 infra 錯誤轉換為 usecase 層語意錯誤
func MapInfraErrorToUseCaseError(err error) error {
	if err == nil {
		return nil
	}

	// 優先比對 CustomError
	switch {
	case errors.Is(err, gateway.ErrMemberNotFound):
		return ErrMemberNotFound
	case errors.Is(err, gateway.ErrMemberAlreadyExists):
		return ErrMemberAlreadyExists
	case errors.Is(err, gateway.ErrMemberUpdateFailed):
		return ErrMemberUpdateFailed
	case errors.Is(err, gateway.ErrMemberDeleteFailed):
		return ErrMemberDeleteFailed
	}

	// fallback 落到底了，真的不認得就回 Unexpected
	return ErrUnexpectedMemberUseCaseFail
}
