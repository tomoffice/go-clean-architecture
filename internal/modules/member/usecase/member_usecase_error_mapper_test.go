package usecase

import (
	"errors"
	"fmt"
	gateway "module-clean/internal/modules/member/interface_adapter/gateway/repository"
	"testing"
)

func TestMapGatewayErrorToUseCaseError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "normal test",
			args: args{
				err: nil,
			},
			wantErr: nil,
		}, {
			name: "ErrGatewayMemberNotFound",
			args: args{
				err: gateway.ErrGatewayMemberNotFound,
			},
			wantErr: ErrUseCaseMemberNotFound,
		}, {
			name: "ErrGatewayMemberAlreadyExists",
			args: args{
				err: gateway.ErrGatewayMemberAlreadyExists,
			},
			wantErr: ErrUseCaseMemberAlreadyExists,
		}, {
			name: "ErrGatewayMemberUpdateFailed",
			args: args{
				err: gateway.ErrGatewayMemberUpdateFailed,
			},
			wantErr: ErrUseCaseMemberUpdateFailed,
		}, {
			name: "ErrGatewayMemberDeleteFailed",
			args: args{
				err: gateway.ErrGatewayMemberDeleteFailed,
			},
			wantErr: ErrUseCaseMemberDeleteFailed,
		}, {
			name: "ErrGatewayMemberDBFailure",
			args: args{
				err: gateway.ErrGatewayMemberDBFailure,
			},
			wantErr: ErrUseCaseMemberDBFailure,
		}, {
			name: "ErrGatewayMemberUnknown",
			args: args{
				err: gateway.ErrGatewayMemberUnknown,
			},
			wantErr: ErrUseCaseMemberUnexpectedFail,
		}, {
			name: "ErrGatewayMemberMappingFailed",
			args: args{
				err: gateway.ErrGatewayMemberMappingFailed,
			},
			wantErr: ErrUseCaseMemberUnexpectedFail,
		}, {
			name: "wrap error",
			args: args{
				err: fmt.Errorf("wrapped: %w 包裹測試", gateway.ErrGatewayMemberNotFound),
			},
			wantErr: ErrUseCaseMemberNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MapGatewayErrorToUseCaseError(tt.args.err); !errors.Is(err, tt.wantErr) {
				t.Errorf("MapGatewayErrorToUseCaseError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
