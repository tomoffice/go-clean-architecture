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

	switch {
	// ─── 業務錯誤 ───────────────────────────────────────────
	case errors.Is(err, gateway.ErrGatewayMemberNotFound):
		return ErrUseCaseMemberNotFound
	case errors.Is(err, gateway.ErrGatewayMemberAlreadyExists):
		return ErrUseCaseMemberAlreadyExists
	case errors.Is(err, gateway.ErrGatewayMemberUpdateFailed):
		return ErrUseCaseMemberUpdateFailed
	case errors.Is(err, gateway.ErrGatewayMemberDeleteFailed):
		return ErrUseCaseMemberDeleteFailed

	// ─── 技術性錯誤：DB 操作異常，歸為非預期 ─────────────
	case errors.Is(err, gateway.ErrGatewayMemberDBError):
		return ErrUseCaseMemberDBError
	}

	// ─── fallback：其他未知錯誤，統一歸為非預期 ────────────
	return ErrUseCaseUnexpectedError
}
