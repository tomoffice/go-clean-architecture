package logger

// Field 代表結構化日誌中的鍵值對欄位。
// 這是純抽象實作，不依賴任何第三方日誌框架，確保可替換性。
//
// Field 用於在日誌訊息中附加額外的結構化資訊，如：
//   - 使用者 ID
//   - 請求 ID
//   - 錯誤代碼
//   - 任何需要記錄的上下文資訊
//
// 設計原則：
// 此類型是抽象層的核心，絕不洩漏底層實作細節。
// 各個 adapter 負責將 logger.Field 轉換為對應框架的欄位格式。
type Field struct {
	Key   string
	Value any
}

// NewField 創建一個新的結構化欄位，支援任意類型的值。
//
// 參數：
//   - key: 欄位的鍵名，應該使用具有描述性的名稱
//   - value: 欄位的值，可以是任意類型（string、int、bool、struct 等）
//
// 返回值：
//   - Field: 可用於日誌記錄的結構化欄位
//
// 使用範例：
//
//	// 基本類型
//	userField := logger.NewField("user_id", "12345")
//	countField := logger.NewField("count", 42)
//	enabledField := logger.NewField("enabled", true)
//
//	// 複雜類型
//	userObj := map[string]interface{}{
//	    "name": "Alice",
//	    "age":  30,
//	}
//	objField := logger.NewField("user", userObj)
//
//	// 使用在日誌中
//	logger.Info("User login", userField, countField)
func NewField(key string, value any) Field {
	return Field{Key: key, Value: value}
}
