package auth

import "errors"

// auth error
var (
	// 初始化階段
	ErrMissingSecretKey  = errors.New("config key is required")
	ErrSecretKeyTooShort = errors.New("config key must be at least 32 characters")

	// 請求階段
	ErrAuthRequired        = errors.New("authorization header required")
	ErrBearerTokenRequired = errors.New("bearer token required")

	// 解析 / 驗證階段
	ErrUnsupportedAlgorithm = errors.New("unsupported signing algorithm") // 不支援的簽名演算法
	ErrTokenExpired         = errors.New("token has expired")             // token 已過期
	ErrSignatureInvalid     = errors.New("signature is invalid")          // 簽名錯誤可能被串改
	ErrParseTokenFailed     = errors.New("failed to parse token")         // 解析 token 時發生錯誤（可能是格式錯誤或其他問題）
	ErrInvalidToken         = errors.New("invalid token")                 // 解析成功但 token 無效（可能是格式錯誤或其他問題）
)

// cors error
