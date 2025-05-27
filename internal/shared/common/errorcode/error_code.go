package errorcode

const (

	// Binding 錯誤（來自 Gin 的 ShouldBindXXX）
	ErrInvalidJSONSyntax = 1000 // JSON 格式錯誤（語法）
	ErrInvalidJSONType   = 1001 // JSON 欄位型別錯誤
	ErrInvalidParams     = 1002 // 其它所有參數綁定錯誤（Query/Form/URI/Header...通用）

	// Validation 錯誤（使用 validator 驗證）
	ErrValidationFailed = 2000 // 欄位驗證失敗（如 required、email 等）
	ErrMissingField     = 2001 // 缺少必要欄位

	// UseCase 業務錯誤
	ErrMemberAlreadyExists         = 3000
	ErrMemberNotFound              = 3001
	ErrMemberUpdateFailed          = 3002
	ErrMemberDeleteFailed          = 3003
	ErrMemberDBFailure             = 3004
	ErrUnexpectedMemberUseCaseFail = 3005

	// 系統錯誤
	ErrInternalServer = 5000
	ErrRequestTimeout = 5001
	ErrContextTimeout = 5002
)
