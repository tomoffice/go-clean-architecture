package auth

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

// allowedAlgs 支援的 JWT 簽名算法
var allowedAlgs = map[string]struct{}{
	"HS256": {},
	"HS384": {},
	"HS512": {},
	"RS256": {},
	"RS384": {},
	"RS512": {},
	"ES256": {},
	"ES384": {},
	"ES512": {},
}

// AuthConfig JWT 認證配置
type AuthConfig struct {
	Secret string `validate:"required,min=32"` // 密鑰，最少 32 字符
}

// AuthMiddleware JWT 認證中間件
type AuthMiddleware[T any] struct {
	config AuthConfig
}

// NewAuthMiddleware 建立 JWT 認證中間件實例
func NewAuthMiddleware[T any](config AuthConfig) (*AuthMiddleware[T], error) {
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, e := range validationErrors {
				if e.Field() == "Secret" {
					if e.Tag() == "required" {
						return nil, ErrMissingSecretKey
					}
					if e.Tag() == "min" {
						return nil, ErrSecretKeyTooShort
					}
				} else {
					return nil, errors.New("validation error: " + e.Error())
				}
			}
		}
		return nil, err
	}
	return &AuthMiddleware[T]{config: config}, nil
}

// HandlerFunc 返回 Gin 中間件函數
func (m *AuthMiddleware[T]) HandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 檢查 Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"error": ErrAuthRequired.Error()})
			return
		}

		// 驗證 Bearer token 格式
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(401, gin.H{"error": ErrBearerTokenRequired.Error()})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 解析並驗證 token
		claims, err := m.decode(tokenString, m.config.Secret)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}

		// 將 claims 存入 context
		ctx.Set("claims", claims)
		ctx.Next()
	}
}

// decode 解析 JWT token 並返回 Claims
func (m *AuthMiddleware[T]) decode(tokenStr string, secret string) (*Claims[T], error) {
	raw := jwt.MapClaims{}

	// 定義密鑰驗證函數
	keyFunc := func(t *jwt.Token) (any, error) {
		if _, ok := allowedAlgs[t.Method.Alg()]; !ok {
			return nil, ErrUnsupportedAlgorithm
		}
		return []byte(secret), nil
	}

	// 解析 token
	token, err := jwt.ParseWithClaims(tokenStr, raw, keyFunc)
	if err != nil {
		// 處理各種錯誤類型
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		if errors.Is(err, ErrUnsupportedAlgorithm) {
			return nil, ErrUnsupportedAlgorithm
		}
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, ErrSignatureInvalid
		}
		return nil, ErrParseTokenFailed
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	// 轉換為結構化的 Claims
	jb, _ := json.Marshal(raw)
	var c Claims[T]
	_ = json.Unmarshal(jb, &c.RegisteredClaims) // 標準 claims
	_ = json.Unmarshal(jb, &c.Private)          // 自定義 claims

	return &c, nil
}
