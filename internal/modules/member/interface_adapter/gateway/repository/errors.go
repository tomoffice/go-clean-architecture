package repository

import "errors"

var (
	ErrMappingFailed       = errors.New("gateway: mapping repo model to entity failed")
	ErrUnknownGateway      = errors.New("gateway: unknown gateway error")
	ErrMemberNotFound      = errors.New("gateway: member not found")
	ErrMemberAlreadyExists = errors.New("gateway: email already exists")
	ErrMemberUpdateFailed  = errors.New("gateway: member update failed")
	ErrMemberDeleteFailed  = errors.New("gateway: member delete failed")
	ErrMemberDBFailure     = errors.New("gateway: member db operation failed")
)
