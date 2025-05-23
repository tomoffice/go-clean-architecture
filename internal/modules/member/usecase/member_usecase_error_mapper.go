package usecase

import (
	"errors"
	gateway "module-clean/internal/modules/member/interface_adapter/gateway/repository"
)

// MapGatewayErrorToUseCaseError 將 Gateway 錯誤轉換為 usecase 層語意錯誤
func MapGatewayErrorToUseCaseError(err error) error {
	if err == nil {
		return nil
	}

	// 優先比對 CustomError
	switch {
	case errors.Is(err, gateway.ErrGatewayMemberNotFound):
		return ErrUseCaseMemberNotFound
	case errors.Is(err, gateway.ErrGatewayMemberAlreadyExists):
		return ErrUseCaseMemberAlreadyExists
	case errors.Is(err, gateway.ErrGatewayMemberUpdateFailed):
		return ErrUseCaseMemberUpdateFailed
	case errors.Is(err, gateway.ErrGatewayMemberDeleteFailed):
		return ErrUseCaseMemberDeleteFailed
	}

	// fallback 落到底了，真的不認得就回 Unexpected
	return ErrUseCaseMemberUnexpectedFail
}
