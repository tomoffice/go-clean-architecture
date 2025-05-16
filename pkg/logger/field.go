package logger

import "go.uber.org/zap"

// Field 結構化欄位的抽象型別（由各 adapter 實作轉換）
type Field interface{}

func String(key, val string) Field {
	return zap.String(key, val)
}

func Int(key string, val int) Field {
	return zap.Int(key, val)
}

func Bool(key string, val bool) Field {
	return zap.Bool(key, val)
}

func Error(err error) Field {
	return zap.Error(err)
}
