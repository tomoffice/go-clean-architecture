package auth

import (
	"strings"
	"testing"
	"time"
)

func TestJWT_GenerateAndParseToken(t *testing.T) {
	// 使用正確的建構函數
	jwt, err := NewWithSecretAndIssuer(
		"this-is-a-very-secure-secret-key-32-chars",
		"test-service",
	)
	if err != nil {
		t.Fatalf("創建 JWT 實例失敗: %v", err)
	}

	// 生成 token
	exp := time.Now().Add(time.Hour)
	token, err := jwt.GenerateToken("user123", exp, map[string]any{
		"role":  "admin",
		"scope": "read:users",
	})
	if err != nil {
		t.Fatalf("生成 token 失敗: %v", err)
	}

	// 驗證 token 格式 (header.payload.signature)
	parts := strings.Split(token.AccessToken, ".")
	if len(parts) != 3 {
		t.Errorf("期望 token 有 3 個部分 (header.payload.signature), 得到 %d 個", len(parts))
	}

	// 解析 token
	claims, customClaims, err := jwt.ParseToken(token.AccessToken)
	if err != nil {
		t.Fatalf("解析 token 失敗: %v", err)
	}

	// 驗證標準 Claims (RFC 7519)
	if claims.Sub != "user123" {
		t.Errorf("期望 sub 為 user123, 得到 %s", claims.Sub)
	}
	if claims.Iss != "test-service" {
		t.Errorf("期望 iss 為 test-service, 得到 %s", claims.Iss)
	}
	if claims.Exp != exp.Unix() {
		t.Errorf("期望 exp 為 %d, 得到 %d", exp.Unix(), claims.Exp)
	}

	// 驗證自定義 Claims
	if customClaims["role"] != "admin" {
		t.Errorf("期望 role 為 admin, 得到 %v", customClaims["role"])
	}
	if customClaims["scope"] != "read:users" {
		t.Errorf("期望 scope 為 read:users, 得到 %v", customClaims["scope"])
	}
}

func TestJWT_ExpiredToken(t *testing.T) {
	jwt, err := NewWithSecretAndIssuer(
		"this-is-a-very-secure-secret-key-32-chars",
		"test-service",
	)
	if err != nil {
		t.Fatalf("創建 JWT 實例失敗: %v", err)
	}

	// 生成已過期的 token
	exp := time.Now().Add(-time.Hour) // 1小時前就過期了
	token, err := jwt.GenerateToken("user123", exp, nil)
	if err != nil {
		t.Fatalf("生成 token 失敗: %v", err)
	}

	// 嘗試解析過期的 token
	_, _, err = jwt.ParseToken(token.AccessToken)
	if err == nil {
		t.Error("期望解析過期 token 時返回錯誤")
	}

	if err != ErrTokenExpired {
		t.Errorf("期望是 token 過期錯誤, 得到: %v", err)
	}

	// 檢查快速過期檢查方法
	if !jwt.IsTokenExpired(token.AccessToken) {
		t.Error("期望 IsTokenExpired 返回 true")
	}
}

func TestJWT_InvalidSignature(t *testing.T) {
	jwt, err := NewWithSecretAndIssuer(
		"this-is-a-very-secure-secret-key-32-chars",
		"test-service",
	)
	if err != nil {
		t.Fatalf("創建 JWT 實例失敗: %v", err)
	}

	// 生成有效 token
	exp := time.Now().Add(time.Hour)
	token, err := jwt.GenerateToken("user123", exp, nil)
	if err != nil {
		t.Fatalf("生成 token 失敗: %v", err)
	}

	// 篡改簽名
	parts := strings.Split(token.AccessToken, ".")
	tamperedToken := parts[0] + "." + parts[1] + "." + "invalid_signature"

	// 嘗試解析被篡改的 token
	_, _, err = jwt.ParseToken(tamperedToken)
	if err == nil {
		t.Error("期望解析被篡改的 token 時返回錯誤")
	}

	if err != ErrInvalidSignature {
		t.Errorf("期望是簽名無效錯誤, 得到: %v", err)
	}
}

func TestJWT_InvalidTokenFormat(t *testing.T) {
	jwt, err := NewWithSecretAndIssuer(
		"this-is-a-very-secure-secret-key-32-chars",
		"test-service",
	)
	if err != nil {
		t.Fatalf("創建 JWT 實例失敗: %v", err)
	}

	invalidTokens := []string{
		"invalid.token",                   // 只有 2 個部分
		"invalid.token.format.extra.part", // 超過 3 個部分
		"",                                // 空字符串
		"invalid",                         // 沒有分隔符
	}

	for _, invalidToken := range invalidTokens {
		_, _, err := jwt.ParseToken(invalidToken)
		if err == nil {
			t.Errorf("期望解析無效 token '%s' 時返回錯誤", invalidToken)
		}
		if err != ErrInvalidToken {
			t.Errorf("期望是無效 token 錯誤, 得到: %v", err)
		}
	}
}

func TestJWT_GenerateTokenWithTTL(t *testing.T) {
	jwt, err := NewWithSecretAndIssuer(
		"this-is-a-very-secure-secret-key-32-chars",
		"test-service",
	)
	if err != nil {
		t.Fatalf("創建 JWT 實例失敗: %v", err)
	}

	// 使用 TTL 生成 token
	ttl := time.Hour
	token, err := jwt.GenerateTokenWithTTL("user123", ttl, map[string]any{"role": "user"})
	if err != nil {
		t.Fatalf("生成 token 失敗: %v", err)
	}

	// 驗證過期時間
	expectedExp := time.Now().Add(ttl).Unix()
	actualExp := token.ExpiresAt.Unix()

	// 允許1秒的誤差
	if abs(expectedExp-actualExp) > 1 {
		t.Errorf("期望過期時間約為 %d, 得到 %d", expectedExp, actualExp)
	}

	// 驗證 token 內容
	claims, customClaims, err := jwt.ParseToken(token.AccessToken)
	if err != nil {
		t.Fatalf("解析 token 失敗: %v", err)
	}

	if claims.Sub != "user123" {
		t.Errorf("期望 sub 為 user123, 得到 %s", claims.Sub)
	}
	if customClaims["role"] != "user" {
		t.Errorf("期望 role 為 user, 得到 %v", customClaims["role"])
	}
}

func TestJWT_RefreshToken(t *testing.T) {
	jwt, err := NewWithSecretAndIssuer(
		"this-is-a-very-secure-secret-key-32-chars",
		"test-service",
	)
	if err != nil {
		t.Fatalf("創建 JWT 實例失敗: %v", err)
	}

	// 生成原始 token
	originalExp := time.Now().Add(time.Hour)
	originalToken, err := jwt.GenerateToken("user123", originalExp, map[string]any{"role": "user"})
	if err != nil {
		t.Fatalf("生成原始 token 失敗: %v", err)
	}

	// 刷新 token
	newExp := time.Now().Add(2 * time.Hour)
	refreshedToken, err := jwt.RefreshToken(originalToken.AccessToken, newExp)
	if err != nil {
		t.Fatalf("刷新 token 失敗: %v", err)
	}

	// 驗證刷新後的 token
	claims, customClaims, err := jwt.ParseToken(refreshedToken.AccessToken)
	if err != nil {
		t.Fatalf("解析刷新後的 token 失敗: %v", err)
	}

	if claims.Sub != "user123" {
		t.Errorf("期望刷新後的 token sub 為 user123, 得到 %s", claims.Sub)
	}
	if customClaims["role"] != "user" {
		t.Errorf("期望刷新後的 token role 為 user, 得到 %v", customClaims["role"])
	}
	if claims.Exp != newExp.Unix() {
		t.Errorf("期望刷新後的 token exp 為 %d, 得到 %d", newExp.Unix(), claims.Exp)
	}
}

func TestNewWithSecret(t *testing.T) {
	secret := "this-is-a-very-secure-secret-key-32-chars"

	jwt, err := NewWithSecret(secret)
	if err != nil {
		t.Fatalf("使用密鑰創建 JWT 實例失敗: %v", err)
	}

	if jwt.config.Secret != secret {
		t.Errorf("期望密鑰為 %s, 得到 %s", secret, jwt.config.Secret)
	}
	if jwt.config.Issuer != "your-service" {
		t.Errorf("期望默認 issuer 為 your-service, 得到 %s", jwt.config.Issuer)
	}
}

func TestJWTConfig_Validation(t *testing.T) {
	// 測試無效配置
	tests := []struct {
		name   string
		config JWTConfig
		valid  bool
	}{
		{
			name: "valid config",
			config: JWTConfig{
				Secret: "this-is-a-very-secure-secret-key-32-chars",
				Issuer: "test-service",
			},
			valid: true,
		},
		{
			name: "secret too short",
			config: JWTConfig{
				Secret: "short",
				Issuer: "test-service",
			},
			valid: false,
		},
		{
			name: "empty issuer",
			config: JWTConfig{
				Secret: "this-is-a-very-secure-secret-key-32-chars",
				Issuer: "",
			},
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New(tt.config)
			if tt.valid && err != nil {
				t.Errorf("期望配置有效，但得到錯誤: %v", err)
			}
			if !tt.valid && err == nil {
				t.Error("期望配置無效，但沒有錯誤")
			}
		})
	}
}

func TestJWTEncoder_EncodeToken(t *testing.T) {
	encoder := NewEncoder("this-is-a-very-secure-secret-key-32-chars", "test-service")

	exp := time.Now().Add(time.Hour).Unix()
	customClaims := map[string]any{"role": "admin"}

	tokenString, err := encoder.EncodeToken("user123", exp, customClaims)
	if err != nil {
		t.Fatalf("編碼 token 失敗: %v", err)
	}

	// 驗證 token 格式
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		t.Errorf("期望 token 有 3 個部分, 得到 %d 個", len(parts))
	}
}

func TestJWTDecoder_DecodeToken(t *testing.T) {
	secret := "this-is-a-very-secure-secret-key-32-chars"
	issuer := "test-service"

	encoder := NewEncoder(secret, issuer)
	decoder := NewDecoder(secret, issuer)

	// 先編碼一個 token
	exp := time.Now().Add(time.Hour).Unix()
	customClaims := map[string]any{"role": "admin", "permissions": []string{"read", "write"}}

	tokenString, err := encoder.EncodeToken("user123", exp, customClaims)
	if err != nil {
		t.Fatalf("編碼 token 失敗: %v", err)
	}

	// 解碼 token
	claims, decodedCustomClaims, err := decoder.DecodeToken(tokenString)
	if err != nil {
		t.Fatalf("解碼 token 失敗: %v", err)
	}

	// 驗證標準 claims
	if claims.Sub != "user123" {
		t.Errorf("期望 sub 為 user123, 得到 %s", claims.Sub)
	}
	if claims.Iss != issuer {
		t.Errorf("期望 iss 為 %s, 得到 %s", issuer, claims.Iss)
	}
	if claims.Exp != exp {
		t.Errorf("期望 exp 為 %d, 得到 %d", exp, claims.Exp)
	}

	// 驗證自定義 claims
	if decodedCustomClaims["role"] != "admin" {
		t.Errorf("期望 role 為 admin, 得到 %v", decodedCustomClaims["role"])
	}
}

func BenchmarkJWT_GenerateToken(b *testing.B) {
	jwt, err := NewWithSecretAndIssuer(
		"this-is-a-very-secure-secret-key-32-chars",
		"bench-service",
	)
	if err != nil {
		b.Fatalf("創建 JWT 實例失敗: %v", err)
	}

	exp := time.Now().Add(time.Hour)
	customClaims := map[string]any{"role": "user"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := jwt.GenerateToken("user123", exp, customClaims)
		if err != nil {
			b.Fatalf("生成 token 失敗: %v", err)
		}
	}
}

func BenchmarkJWT_ParseToken(b *testing.B) {
	jwt, err := NewWithSecretAndIssuer(
		"this-is-a-very-secure-secret-key-32-chars",
		"bench-service",
	)
	if err != nil {
		b.Fatalf("創建 JWT 實例失敗: %v", err)
	}

	exp := time.Now().Add(time.Hour)
	token, err := jwt.GenerateToken("user123", exp, map[string]any{"role": "user"})
	if err != nil {
		b.Fatalf("生成 token 失敗: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := jwt.ParseToken(token.AccessToken)
		if err != nil {
			b.Fatalf("解析 token 失敗: %v", err)
		}
	}
}

// 輔助函數
func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}
