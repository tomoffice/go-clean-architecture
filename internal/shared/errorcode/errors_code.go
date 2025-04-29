package errorcode

const (
	// 輸入驗證錯誤
	ErrInvalidJSONSyntax   = 1000
	ErrInvalidJSONType     = 1001
	ErrInvalidJSONInput    = 1002
	ErrValidationFailed    = 1003
	ErrMissingField        = 1004
	ErrInvalidQueryParams  = 1005
	ErrInvalidURIParams    = 1006
	ErrInvalidFormData     = 1007
	ErrInvalidHeaderParams = 1008

	// usecase錯誤
	ErrMemberAlreadyExists         = 2000
	ErrMemberNotFound              = 2001
	ErrMemberUpdateFailed          = 2002
	ErrMemberDeleteFailed          = 2003
	ErrMemberDBFailure             = 2004
	ErrUnexpectedMemberUseCaseFail = 2005
	// 例外捕獲
	ErrInternalServer = 5000

	// IO錯誤
	ErrRequestTimeout = 5001
	ErrContextTimeout = 5002
)
