package http

import (
	"errors"
	"module-clean/internal/modules/member/usecase"
	"module-clean/internal/shared/common/errorcode"
	sharederrors "module-clean/internal/shared/common/errordefs"
)

func MapMemberUseCaseError(err error) (int, string) {
	switch {
	case errors.Is(err, usecase.ErrUseCaseMemberNotFound):
		return errorcode.ErrMemberNotFound, usecase.ErrUseCaseMemberNotFound.Error()
	case errors.Is(err, usecase.ErrUseCaseMemberAlreadyExists):
		return errorcode.ErrMemberAlreadyExists, usecase.ErrUseCaseMemberAlreadyExists.Error()
	case errors.Is(err, usecase.ErrUseCaseMemberUpdateFailed):
		return errorcode.ErrMemberUpdateFailed, usecase.ErrUseCaseMemberUpdateFailed.Error()
	case errors.Is(err, usecase.ErrUseCaseMemberDeleteFailed):
		return errorcode.ErrMemberDeleteFailed, usecase.ErrUseCaseMemberDeleteFailed.Error()
	case errors.Is(err, usecase.ErrUseCaseMemberUnexpectedFail):
		return errorcode.ErrUnexpectedMemberUseCaseFail, usecase.ErrUseCaseMemberUnexpectedFail.Error()
	case errors.Is(err, usecase.ErrUseCaseMemberDBFailure):
		return errorcode.ErrMemberDBFailure, usecase.ErrUseCaseMemberDBFailure.Error()
	default:
		return errorcode.ErrInternalServer, sharederrors.ErrInternalServer.Error()
	}
}
