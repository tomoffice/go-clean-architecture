package errorcode

const (

	// Binding 錯誤（來自 Gin 的 ShouldBindXXX）
	ErrInvalidJSONSyntax   = 1000 // JSON 格式錯誤（語法）
	ErrInvalidJSONType     = 1001 // JSON 欄位型別錯誤
	ErrInvalidJSONInput    = 1002 // JSON 綁定失敗（無法解析）
	ErrInvalidQueryParams  = 1003 // Query 綁定錯誤
	ErrInvalidFormData     = 1004 // Form 綁定錯誤
	ErrInvalidURIParams    = 1005 // URI 綁定錯誤
	ErrInvalidHeaderParams = 1006 // Header 綁定錯誤
	ErrUnexpectedBinding   = 1009 // Gin Binding 無法分類錯誤

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
