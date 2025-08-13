package seq

import (
	"go.uber.org/zap"

	"github.com/tomoffice/go-clean-architecture/pkg/logger"
)

// toZapFields 將抽象的 logger.Field 切片轉換為 zap.Field 切片
// 這確保了抽象層不洩漏 zap 實作細節
func toZapFields(fields ...logger.Field) []zap.Field {
	if len(fields) == 0 {
		return nil
	}
	
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}