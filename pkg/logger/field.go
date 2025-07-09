package logger

import (
	"go.uber.org/zap"
)

// Field 結構化欄位類型
// 目前基於 zap.Field，如果未來需要支援其他日誌框架可以改為介面
type Field = zap.Field

func NewField[T any](key string, value T) Field {
	return zap.Any(key, value)
}

// 如果未來需要支援其他日誌框架，可以改為：
// type Field interface {
//     // 定義必要的方法
// }
//
// 然後各個 adapter 實作對應的轉換邏輯
