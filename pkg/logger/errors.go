package logger

import "errors"

// 定義 logger 套件的錯誤類型
var (
	// ErrNoValidLoggers 當沒有有效的 logger 可用時回傳
	ErrNoValidLoggers = errors.New("沒有有效的 logger 可用")

	// ErrInvalidConfig 當設定參數無效時回傳
	ErrInvalidConfig = errors.New("無效的 logger 設定")

	// ErrLoggerCreationFailed 當建立 logger 失敗時回傳
	ErrLoggerCreationFailed = errors.New("建立 logger 失敗")

	// ErrUnsupportedFormat 當不支援的格式時回傳
	ErrUnsupportedFormat = errors.New("不支援的日誌格式")
)