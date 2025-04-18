package errorcode

const (
	// 通用
	ErrInvalidJSONInput = 1000
	ErrValidationFailed = 1001
	ErrInternalServer   = 1500

	// 前端可用的業務錯誤碼
	ErrEmailExists        = 2001
	ErrMemberNotFound     = 2002
	ErrMemberUpdateFailed = 2003
	ErrMemberDeleteFailed = 2004

	// UseCase 運行錯誤
	ErrTokenInvalid = 3001
)
