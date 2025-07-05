package auth

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"strings"
)

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

type AuthConfig struct {
	Secret string `validate:"required,min=32"` // JWT 密鑰，至少 32 字符長
}

// AuthMiddleware JWT 認證的中間件
//   - config: 認證配置，包括密鑰和其他選項
type AuthMiddleware[T any] struct {
	config AuthConfig
}

// NewAuthMiddleware 建立 JWT 認證中間件
func NewAuthMiddleware[T any](config AuthConfig) (*AuthMiddleware[T], error) {
	var validate = validator.New()
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

func (m *AuthMiddleware[T]) HandlerFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"error": ErrAuthRequired.Error()})
			return
		}

		// 處理 Bearer 前綴
		if !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.AbortWithStatusJSON(401, gin.H{"error": ErrBearerTokenRequired.Error()})
			return
		}
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// 解析 token
		claims, err := m.decode(tokenString, m.config.Secret)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": err.Error()})
			return
		}
		// 將解析出的聲明保存到上下文中
		ctx.Set("claims", claims)
		ctx.Next()
	}
}

// decode 解析 JWT token 並返回Claims
//   - tokenStr: 要解析的 JWT token 字符串
//   - secret: 用於驗證簽名的密鑰
func (m *AuthMiddleware[T]) decode(tokenStr string, secret string) (*Claims[T], error) {
	raw := jwt.MapClaims{}
	keyFunc := func(t *jwt.Token) (any, error) {
		if _, ok := allowedAlgs[t.Method.Alg()]; !ok {
			return nil, ErrUnsupportedAlgorithm
		}
		return []byte(secret), nil
	}
	token, err := jwt.ParseWithClaims(tokenStr, raw, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired // token 已過期
		}
		if errors.Is(err, ErrUnsupportedAlgorithm) {
			return nil, ErrUnsupportedAlgorithm // 算法不匹配
		}
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, ErrSignatureInvalid // 簽章驗證失敗（被串改）
		}
		return nil, ErrParseTokenFailed // 其他解析錯誤
	}
	if !token.Valid {
		return nil, ErrInvalidToken // 解析沒錯但 token 無效
	}
	jb, _ := json.Marshal(raw)

	// 解析聲明
	var c Claims[T]
	_ = json.Unmarshal(jb, &c.RegisteredClaims) // 標準聲明
	_ = json.Unmarshal(jb, &c.Private)          // 自定義聲明

	return &c, nil
}
