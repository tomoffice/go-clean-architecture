package logger

import (
	"go.uber.org/zap"
)

// Field 代表結構化日誌中的鍵值對欄位。
// 目前基於 zap.Field 實作，提供高效能的結構化欄位支援。
//
// Field 用於在日誌訊息中附加額外的結構化資訊，如：
//   - 使用者 ID
//   - 請求 ID
//   - 錯誤代碼
//   - 任何需要記錄的上下文資訊
//
// 設計考量：
// 目前直接使用 zap.Field 以獲得最佳效能，如果未來需要支援其他日誌框架，
// 可以將此類型改為介面，並在各個 adapter 中實作相應的轉換邏輯。
type Field = zap.Field

// NewField 創建一個新的結構化欄位，支援任意類型的值。
// 此函數使用 Go 泛型確保型別安全，並自動處理值的序列化。
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
func NewField[T any](key string, value T) Field {
	return zap.Any(key, value)
}

// 設計備註：
// 如果未來需要支援其他日誌框架，可以將 Field 改為介面：
//
// type Field interface {
//     Key() string
//     Value() interface{}
//     Type() FieldType
// }
//
// 然後在各個 adapter 中實作對應的轉換邏輯，
// 以支援不同日誌框架的欄位格式。