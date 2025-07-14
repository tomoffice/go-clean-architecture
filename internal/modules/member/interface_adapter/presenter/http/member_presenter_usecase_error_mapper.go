package http

import (
	"errors"
	"github.com/tomoffice/go-clean-architecture/internal/modules/member/usecase"
	"github.com/tomoffice/go-clean-architecture/internal/shared/errorcode"
	sharederrors "github.com/tomoffice/go-clean-architecture/internal/shared/errordefs"
)

func MapMemberUseCaseToPresenterError(err error) (int, string) {
	switch {
	case errors.Is(err, usecase.ErrMemberNotFound):
		return errorcode.ErrMemberNotFound, usecase.ErrMemberNotFound.Error()
	case errors.Is(err, usecase.ErrMemberAlreadyExists):
		return errorcode.ErrMemberAlreadyExists, usecase.ErrMemberAlreadyExists.Error()
	case errors.Is(err, usecase.ErrMemberNoEffect):
		return errorcode.ErrMemberNoEffect, usecase.ErrMemberNoEffect.Error()
	case errors.Is(err, usecase.ErrMemberDBError):
		return errorcode.ErrMemberDBError, usecase.ErrMemberDBError.Error()
	case errors.Is(err, usecase.ErrMemberUnexpectedError):
		return errorcode.ErrUnexpectedMemberUseCaseError, usecase.ErrMemberUnexpectedError.Error()
	case errors.Is(err, usecase.ErrMemberUpdateSameEmail):
		return errorcode.ErrMemberUpdateSameEmail, usecase.ErrMemberUpdateSameEmail.Error()
	case errors.Is(err, usecase.ErrMemberEmailAlreadyExists):
		return errorcode.ErrMemberEmailAlreadyExists, usecase.ErrMemberEmailAlreadyExists.Error()
	case errors.Is(err, usecase.ErrMemberPasswordIncorrect):
		return errorcode.ErrMemberPasswordIncorrect, usecase.ErrMemberPasswordIncorrect.Error()
	default:
		return errorcode.ErrInternalServer, sharederrors.ErrInternalServer.Error()
	}
}
