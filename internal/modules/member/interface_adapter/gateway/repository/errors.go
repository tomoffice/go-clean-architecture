package repository

import "errors"

var (
	ErrGatewayMemberMappingFailed = errors.New("gateway: mapping repo model to entity failed")
	ErrGatewayMemberUnknown       = errors.New("gateway: unknown gateway error")
	ErrGatewayMemberNotFound      = errors.New("gateway: member not found")
	ErrGatewayMemberAlreadyExists = errors.New("gateway: email already exists")
	ErrGatewayMemberUpdateFailed  = errors.New("gateway: member update failed")
	ErrGatewayMemberDeleteFailed  = errors.New("gateway: member delete failed")
	ErrGatewayMemberDBFailure     = errors.New("gateway: member db operation failed")
)
