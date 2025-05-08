package http

import (
	"errors"
	"module-clean/internal/modules/member/usecase"
	"module-clean/internal/shared/common/errorcode"
	sharederrors "module-clean/internal/shared/common/errordefs"
)

func MapMemberUseCaseError(err error) (int, string) {
	switch {
	case errors.Is(err, usecase.ErrMemberNotFound):
		return errorcode.ErrMemberNotFound, usecase.ErrMemberNotFound.Error()
	case errors.Is(err, usecase.ErrMemberAlreadyExists):
		return errorcode.ErrMemberAlreadyExists, usecase.ErrMemberAlreadyExists.Error()
	case errors.Is(err, usecase.ErrMemberUpdateFailed):
		return errorcode.ErrMemberUpdateFailed, usecase.ErrMemberUpdateFailed.Error()
	case errors.Is(err, usecase.ErrMemberDeleteFailed):
		return errorcode.ErrMemberDeleteFailed, usecase.ErrMemberDeleteFailed.Error()
	case errors.Is(err, usecase.ErrUnexpectedMemberUseCaseFail):
		return errorcode.ErrUnexpectedMemberUseCaseFail, usecase.ErrUnexpectedMemberUseCaseFail.Error()
	case errors.Is(err, usecase.ErrMemberDBFailure):
		return errorcode.ErrMemberDBFailure, usecase.ErrMemberDBFailure.Error()
	default:
		return errorcode.ErrInternalServer, sharederrors.ErrInternalServer.Error()
	}
}
