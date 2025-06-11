package http

import (
	"errors"
	"module-clean/internal/modules/member/usecase"
	"module-clean/internal/shared/common/errorcode"
	sharederrors "module-clean/internal/shared/common/errordefs"
)

func MapMemberUseCaseToPresenterError(err error) (int, string) {
	switch {
	case errors.Is(err, usecase.ErrUseCaseMemberNotFound):
		return errorcode.ErrMemberNotFound, usecase.ErrUseCaseMemberNotFound.Error()
	case errors.Is(err, usecase.ErrUseCaseMemberAlreadyExists):
		return errorcode.ErrMemberAlreadyExists, usecase.ErrUseCaseMemberAlreadyExists.Error()
	case errors.Is(err, usecase.ErrUseCaseMemberNoEffect):
		return errorcode.ErrMemberNoEffect, usecase.ErrUseCaseMemberNoEffect.Error()
	case errors.Is(err, usecase.ErrUseCaseMemberDBError):
		return errorcode.ErrMemberDBError, usecase.ErrUseCaseMemberDBError.Error()
	case errors.Is(err, usecase.ErrUseCaseMemberUnexpectedError):
		return errorcode.ErrUnexpectedMemberUseCaseError, usecase.ErrUseCaseMemberUnexpectedError.Error()
	case errors.Is(err, usecase.ErrUseCaseMemberGatewayError):
		return errorcode.ErrMemberGatewayError, usecase.ErrUseCaseMemberGatewayError.Error()
	case errors.Is(err, usecase.ErrUseCaseMemberUpdateSameEmail):
		return errorcode.ErrMemberUpdateSameEmail, usecase.ErrUseCaseMemberUpdateSameEmail.Error()
	case errors.Is(err, usecase.ErrUseCaseMemberEmailAlreadyExists):
		return errorcode.ErrMemberEmailAlreadyExists, usecase.ErrUseCaseMemberEmailAlreadyExists.Error()
	case errors.Is(err, usecase.ErrUseCaseMemberPasswordIncorrect):
		return errorcode.ErrMemberPasswordIncorrect, usecase.ErrUseCaseMemberPasswordIncorrect.Error()
	default:
		return errorcode.ErrInternalServer, sharederrors.ErrInternalServer.Error()
	}
}
