package controller

import (
	"github.com/stretchr/testify/assert"
	"module-clean/internal/shared/common/errorcode"
	"net/http"
	"testing"
)

func TestMapErrorCodeToHTTPStatus(t *testing.T) {
	type args struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Unknown Error Code",
			args: args{
				code: 999999,
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Binding Error",
			args: args{
				code: errorcode.ErrInvalidJSONType,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Validation Error",
			args: args{
				code: errorcode.ErrInvalidParams,
			},
			want: http.StatusBadRequest,
		},
		{
			name: "UseCase Error - Not Found",
			args: args{
				code: errorcode.ErrMemberNotFound,
			},
			want: http.StatusNotFound,
		},
		{
			name: "UseCase Error - Already Exists",
			args: args{
				code: errorcode.ErrMemberAlreadyExists,
			},
			want: http.StatusConflict,
		},
		{
			name: "UseCase Error - No Effect",
			args: args{
				code: errorcode.ErrMemberNoEffect,
			},
			want: http.StatusUnprocessableEntity,
		},
		{
			name: "UseCase Error - Email Already Exists",
			args: args{
				code: errorcode.ErrMemberEmailAlreadyExists,
			},
			want: http.StatusConflict,
		},
		{
			name: "UseCase Error - Update Same Email",
			args: args{
				code: errorcode.ErrMemberUpdateSameEmail,
			},
			want: http.StatusConflict,
		},
		{
			name: "UseCase Error - Password Incorrect",
			args: args{
				code: errorcode.ErrMemberPasswordIncorrect,
			},
			want: http.StatusUnauthorized,
		},
		{
			name: "UseCase Error - Update Same Password",
			args: args{
				code: errorcode.ErrMemberUpdateSamePassword,
			},
			want: http.StatusConflict,
		},
		{
			name: "UseCase Error - Unexpected Error",
			args: args{
				code: errorcode.ErrUnexpectedMemberUseCaseError,
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "System Error - Request Timeout",
			args: args{
				code: errorcode.ErrRequestTimeout,
			},
			want: http.StatusGatewayTimeout,
		},
		{
			name: "System Error - Context Timeout",
			args: args{
				code: errorcode.ErrContextTimeout,
			},
			want: http.StatusGatewayTimeout,
		},
		{
			name: "System Error - Internal Server Error",
			args: args{
				code: errorcode.ErrInternalServer,
			},
			want: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, MapErrorCodeToHTTPStatus(tt.args.code), "MapErrorCodeToHTTPStatus(%v)", tt.args.code)
		})
	}
}
