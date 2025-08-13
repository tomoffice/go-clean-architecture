package seq

import (
	"go.uber.org/zap/zapcore"
	"math"
	"time"
)

// extractValue 從 zapcore.Field 中擷取出對應的 Go 原生值
// 在 core_adapter.Write 裡會用這個結果去填入 logrus.Fields
func extractValue(f zapcore.Field) any {
	//fmt.Printf("Field Type: %v, Key: %s, Interface: %T\n", f.Type, f.Key, f.Interface)
	switch f.Type {
	case zapcore.StringType:
		return f.String
	case zapcore.Int64Type, zapcore.Int32Type, zapcore.Int16Type, zapcore.Int8Type:
		return f.Integer
	case zapcore.BoolType:
		return f.Integer != 0
	case zapcore.Float64Type:
		return math.Float64frombits(uint64(f.Integer))
	case zapcore.DurationType:
		return time.Duration(f.Integer)
	case zapcore.TimeType:
		return time.Unix(0, f.Integer).Format(time.RFC3339Nano)
	case zapcore.ErrorType:
		if f.Interface != nil {
			if err, ok := f.Interface.(error); ok {
				return err.Error()
			}
		}
		return f.String
	case zapcore.ReflectType, zapcore.ObjectMarshalerType, zapcore.NamespaceType:
		return f.Interface
	default:
		return nil
	}
}