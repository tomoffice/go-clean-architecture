package errorcode

const (
	// 輸入驗證錯誤
	ErrInvalidJSONSyntax  = 1000 // JSON 格式錯誤
	ErrInvalidJSONType    = 1001 // JSON 欄位型別錯誤
	ErrInvalidJSONInput   = 1002 // JSON 輸入錯誤 fallback
	ErrValidationFailed   = 1003
	ErrMissingField       = 1004
	ErrInvalidQueryParams = 1005

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
