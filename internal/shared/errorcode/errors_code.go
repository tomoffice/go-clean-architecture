package errorcode

const (
	// 輸入驗證錯誤
	ErrInvalidJSONSyntax  = 1000 // JSON 格式錯誤
	ErrInvalidJSONType    = 1001 // JSON 欄位型別錯誤
	ErrInvalidJSONInput   = 1002 // JSON 輸入錯誤 fallback
	ErrValidationFailed   = 1003
	ErrMissingField       = 1004
	ErrInvalidQueryParams = 1005

	// 業務邏輯錯誤
	ErrMemberAlreadyExists = 2001
	ErrMemberNotFound      = 2002
	ErrMemberUpdateFailed  = 2003
	ErrMemberDeleteFailed  = 2004

	// 例外捕獲
	ErrInternalServer = 5000

	// IO錯誤
	ErrRequestTimeout = 5001
	ErrContextTimeout = 5002
)
