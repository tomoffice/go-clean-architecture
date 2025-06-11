package errorcode

// Binding 錯誤（來自 Gin 的 ShouldBindXXX）
const (
	ErrInvalidJSONSyntax = 1000 // JSON 格式錯誤（語法）
	ErrInvalidJSONType   = 1001 // JSON 欄位型別錯誤
	ErrInvalidParams     = 1002 // 其它所有參數綁定錯誤（Query/Form/URI/Header...通用）
)

// Validation 錯誤（使用 validator 驗證）
const (
	ErrValidationFailed = 2000 // 欄位驗證失敗（如 required、email 等）
	ErrMissingField     = 2001 // 缺少必要欄位
)

// UseCase 層相關業務錯誤
const (
	ErrMemberNotFound               = 3000 // 會員不存在
	ErrMemberAlreadyExists          = 3001 // 會員已存在
	ErrMemberNoEffect               = 3002 // 更新/刪除無影響
	ErrMemberDBError                = 3003 // DB 錯誤
	ErrMemberGatewayError           = 3005 // Gateway 層錯誤
	ErrUnexpectedMemberUseCaseError = 3006 // 非預期 UseCase 錯誤
	ErrMemberUpdateSameEmail        = 3007 // 嘗試更新為同一 Email
	ErrMemberEmailAlreadyExists     = 3008 // Email 已被佔用
	ErrMemberPasswordIncorrect      = 3010 // 密碼錯誤
	ErrMemberUpdateSamePassword     = 3009 // 嘗試更新為同一密碼
)

// 系統錯誤
const (
	ErrInternalServer = 5000 // 系統內部錯誤
	ErrRequestTimeout = 5001 // 請求逾時
	ErrContextTimeout = 5002 // context 超時
)
