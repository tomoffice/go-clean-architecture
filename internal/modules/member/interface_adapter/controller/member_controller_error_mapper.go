package controller

import (
	"module-clean/internal/shared/errorcode"
	"net/http"
)

func MapErrorCodeToHTTPStatus(code int) int {
	switch {
	// Binding
	case code >= 1000 && code < 2000:
		return http.StatusBadRequest

	// Validation → 400
	case code >= 2000 && code < 3000:
		return http.StatusBadRequest

	// UseCase → 404, 409, or 500
	case code == errorcode.ErrMemberNotFound:
		return http.StatusNotFound
	case code == errorcode.ErrMemberAlreadyExists:
		return http.StatusConflict
	case code == errorcode.ErrMemberNoEffect:
		return http.StatusUnprocessableEntity
	case code == errorcode.ErrMemberEmailAlreadyExists:
		return http.StatusConflict
	case code == errorcode.ErrMemberUpdateSameEmail:
		return http.StatusConflict
	case code == errorcode.ErrMemberPasswordIncorrect:
		return http.StatusUnauthorized
	case code == errorcode.ErrMemberUpdateSamePassword:
		return http.StatusConflict
	case code >= 3000 && code < 4000:
		return http.StatusInternalServerError

	// 系統錯誤 → 500 or 504
	case code == errorcode.ErrRequestTimeout || code == errorcode.ErrContextTimeout:
		return http.StatusGatewayTimeout
	case code >= 5000 && code < 6000:
		return http.StatusInternalServerError

	// fallback
	default:
		return http.StatusInternalServerError
	}
}
