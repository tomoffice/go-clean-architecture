package domain

import "errors"

var ( // 語意化的會員相關錯誤
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrMemberNotFound     = errors.New("member not found")
	ErrMemberUpdateFailed = errors.New("member update failed")
	ErrMemberDeleteFailed = errors.New("member delete failed")
)
