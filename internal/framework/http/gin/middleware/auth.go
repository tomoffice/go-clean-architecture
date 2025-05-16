package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthType 定義認證類型
type AuthType string

const (
	AuthTypeBearer AuthType = "Bearer"
	AuthTypeAPIKey AuthType = "ApiKey"
	AuthTypeBasic  AuthType = "Basic"
)

// AuthMiddleware 實現認證的 middleware
// 支援 JWT、API Key 等多種認證方式
type AuthMiddleware struct {
	authService AuthService
	authType    AuthType
	skipPaths   []string

	// JWT 相關配置
	secretKey     []byte
	tokenLookup   string // header:Authorization, query:token, cookie:jwt
	tokenHeadName string // Bearer
	timeFunc      func() time.Time
}

// AuthConfig 認證 middleware 配置
type AuthConfig struct {
	AuthService AuthService
	AuthType    AuthType
	SkipPaths   []string

	// JWT 相關配置
	SecretKey     string
	TokenLookup   string // 預設: "header:Authorization"
	TokenHeadName string // 預設: "Bearer"
	TimeFunc      func() time.Time
}

// JWTClaims JWT claims 結構
type JWTClaims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

// DefaultAuthConfig 提供預設的認證配置
func DefaultAuthConfig(authService AuthService) AuthConfig {
	return AuthConfig{
		AuthService:   authService,
		AuthType:      AuthTypeBearer,
		TokenLookup:   "header:Authorization",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
		SkipPaths:     []string{},
	}
}

// NewAuthMiddleware 建立認證 middleware
func NewAuthMiddleware(config AuthConfig) Middleware {
	if config.TokenLookup == "" {
		config.TokenLookup = "header:Authorization"
	}
	if config.TokenHeadName == "" {
		config.TokenHeadName = "Bearer"
	}
	if config.TimeFunc == nil {
		config.TimeFunc = time.Now
	}

	return &AuthMiddleware{
		authService:   config.AuthService,
		authType:      config.AuthType,
		skipPaths:     config.SkipPaths,
		secretKey:     []byte(config.SecretKey),
		tokenLookup:   config.TokenLookup,
		tokenHeadName: config.TokenHeadName,
		timeFunc:      config.TimeFunc,
	}
}

// Handle 實現 Middleware 介面
func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 檢查是否需要跳過認證
		if m.shouldSkip(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 根據認證類型進行驗證
		var userID string
		var err error

		switch m.authType {
		case AuthTypeBearer:
			userID, err = m.handleJWTAuth(c)
		case AuthTypeAPIKey:
			userID, err = m.handleAPIKeyAuth(c)
		case AuthTypeBasic:
			userID, err = m.handleBasicAuth(c)
		default:
			err = errors.New("不支援的認證類型")
		}

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": err.Error(),
			})
			c.Abort()
			return
		}

		// 將使用者資訊存入 Gin Context
		c.Set("user_id", userID)
		c.Set("authenticated", true)

		c.Next()
	}
}

// handleJWTAuth 處理 JWT 認證
func (m *AuthMiddleware) handleJWTAuth(c *gin.Context) (string, error) {
	// 提取 token
	token, err := m.extractToken(c)
	if err != nil {
		return "", errors.New("未提供認證令牌")
	}

	// 驗證 JWT token
	claims, err := m.validateJWT(token)
	if err != nil {
		return "", errors.New("無效的認證令牌: " + err.Error())
	}

	// 將 claims 資訊存入 context
	c.Set("jwt_claims", claims)
	c.Set("username", claims.Username)
	c.Set("email", claims.Email)
	c.Set("roles", claims.Roles)

	return claims.UserID, nil
}

// handleAPIKeyAuth 處理 API Key 認證
func (m *AuthMiddleware) handleAPIKeyAuth(c *gin.Context) (string, error) {
	apiKey := c.GetHeader("X-API-Key")
	if apiKey == "" {
		return "", errors.New("未提供 API Key")
	}

	userID, err := m.authService.ValidateAPIKey(apiKey)
	if err != nil {
		return "", errors.New("無效的 API Key")
	}

	return userID, nil
}

// handleBasicAuth 處理基本認證
func (m *AuthMiddleware) handleBasicAuth(c *gin.Context) (string, error) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		return "", errors.New("未提供基本認證資訊")
	}

	userID, err := m.authService.ValidateBasicAuth(username, password)
	if err != nil {
		return "", errors.New("認證失敗")
	}

	return userID, nil
}

// extractToken 從請求中提取 token
func (m *AuthMiddleware) extractToken(c *gin.Context) (string, error) {
	parts := strings.Split(m.tokenLookup, ":")
	if len(parts) != 2 {
		return "", errors.New("無效的 token 查找配置")
	}

	key := parts[0]
	val := parts[1]

	tokenString := ""

	switch key {
	case "header":
		tokenString = c.Request.Header.Get(val)
	case "query":
		tokenString = c.Query(val)
	case "cookie":
		cookie, err := c.Cookie(val)
		if err != nil {
			return "", err
		}
		tokenString = cookie
	default:
		return "", errors.New("無效的 token 查找方式")
	}

	if tokenString == "" {
		return "", errors.New("token 不存在")
	}

	// 如果是從 header 取得，移除 Bearer 前綴
	if key == "header" && m.tokenHeadName != "" {
		parts := strings.SplitN(tokenString, " ", 2)
		if len(parts) != 2 || parts[0] != m.tokenHeadName {
			return "", errors.New("無效的 token 格式")
		}
		tokenString = parts[1]
	}

	return tokenString, nil
}

// validateJWT 驗證 JWT token
func (m *AuthMiddleware) validateJWT(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 驗證簽名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("無效的簽名方法")
		}
		return m.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("無效的 token")
}

// shouldSkip 檢查路徑是否需要跳過認證
func (m *AuthMiddleware) shouldSkip(path string) bool {
	for _, skipPath := range m.skipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

// --- 輔助函數 ---

// GetUserID 從 Gin Context 獲取使用者 ID
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}

	id, ok := userID.(string)
	return id, ok
}

// GetJWTClaims 從 Gin Context 獲取 JWT claims
func GetJWTClaims(c *gin.Context) (*JWTClaims, bool) {
	claims, exists := c.Get("jwt_claims")
	if !exists {
		return nil, false
	}

	jwtClaims, ok := claims.(*JWTClaims)
	return jwtClaims, ok
}

// IsAuthenticated 檢查請求是否已認證
func IsAuthenticated(c *gin.Context) bool {
	authenticated, exists := c.Get("authenticated")
	if !exists {
		return false
	}

	auth, ok := authenticated.(bool)
	return ok && auth
}

// GenerateJWT 生成 JWT token
func GenerateJWT(claims *JWTClaims, secretKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secretKey))
}

// GenerateJWTWithExpiry 生成有效期限的 JWT token
func GenerateJWTWithExpiry(userID, username, email string, roles []string, secretKey string, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	return GenerateJWT(claims, secretKey)
}
