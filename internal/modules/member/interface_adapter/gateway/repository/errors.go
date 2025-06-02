package repository

import "errors"

var (
	ErrGatewayMemberMappingError    = errors.New("gateway: mapping repo model to entity failed")
	ErrGatewayMemberUnexpectedError = errors.New("gateway: member gateway unexpected error")
	ErrGatewayMemberNotFound        = errors.New("gateway: member not found")
	ErrGatewayMemberAlreadyExists   = errors.New("gateway: member already exists")
	ErrGatewayMemberUpdateFailed    = errors.New("gateway: member update failed")
	ErrGatewayMemberDeleteFailed    = errors.New("gateway: member delete failed")
	ErrGatewayMemberDBError         = errors.New("gateway: member db operation error")
)
