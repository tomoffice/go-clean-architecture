package http

import (
	"errors"
	"module-clean/internal/member/domain"
	"module-clean/internal/shared/errorcode"
)

func MapErrorToCode(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrMemberNotFound):
		return errorcode.ErrMemberNotFound, domain.ErrMemberNotFound.Error()
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		return errorcode.ErrEmailExists, domain.ErrEmailAlreadyExists.Error()
	case errors.Is(err, domain.ErrMemberUpdateFailed):
		return errorcode.ErrMemberUpdateFailed, domain.ErrMemberUpdateFailed.Error()
	case errors.Is(err, domain.ErrMemberDeleteFailed):
		return errorcode.ErrMemberDeleteFailed, domain.ErrMemberDeleteFailed.Error()
	default:
		return errorcode.ErrInternalServer, errors.New("internal server error").Error()
	}
}
