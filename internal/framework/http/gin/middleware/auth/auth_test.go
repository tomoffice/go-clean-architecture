package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandlerFunc(t *testing.T) {
	type privateClaims struct {
		AA string `json:"aa"`
		BB string `json:"bb"`
		CC string `json:"cc"`
	}
	longSecret := strings.Repeat("x", 32)
	tests := []struct {
		name           string
		config         AuthConfig
		setupHeader    func(req *http.Request)
		wantHTTPStatus int
		wantJSON       string
		wantClaims     *Claims[privateClaims]
	}{
		{
			name:   "normal case",
			config: AuthConfig{Secret: longSecret},
			setupHeader: func(req *http.Request) {
				cs := jwt.MapClaims{
					"iss": "test",
					"aud": "single",
					"exp": time.Now().Add(1 * time.Second).Unix(), // 過期時間
				}
				token, _ := genToken(t, longSecret, cs, jwt.SigningMethodHS256)
				req.Header.Set("Authorization", "Bearer "+token)
			},
			wantHTTPStatus: 200,
			wantJSON:       "",
			wantClaims: &Claims[privateClaims]{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "test",
					Audience:  jwt.ClaimStrings{"single"},
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Second)),
				},
			},
		},
		{
			name:   "no authorization header",
			config: AuthConfig{Secret: longSecret},
			setupHeader: func(req *http.Request) {
				req.Header.Set("", "")
			},
			wantHTTPStatus: http.StatusUnauthorized,
			wantJSON:       fmt.Sprintf(`{"error":"%s"}`, ErrAuthRequired),
			wantClaims:     nil,
		},
		{
			name:   "missing Bearer prefix",
			config: AuthConfig{Secret: longSecret},
			setupHeader: func(req *http.Request) {
				req.Header.Set("Authorization", "test ")
			},
			wantHTTPStatus: http.StatusUnauthorized,
			wantJSON:       fmt.Sprintf(`{"error":"%s"}`, ErrBearerTokenRequired),
			wantClaims:     nil,
		},
		{
			name:   "invalid token format",
			config: AuthConfig{Secret: longSecret},
			setupHeader: func(req *http.Request) {
				req.Header.Set("Authorization", "Bearer invalid.token.format")
			},
			wantHTTPStatus: http.StatusUnauthorized,
			wantJSON:       fmt.Sprintf(`{"error":"%s"}`, ErrParseTokenFailed),
			wantClaims:     nil,
		},
		{
			name:   "invalid token - expired",
			config: AuthConfig{Secret: longSecret},
			setupHeader: func(req *http.Request) {
				cs := jwt.MapClaims{
					"iss": "test",
					"aud": "single",
					"exp": time.Now().Add(-1 * time.Second).Unix(), // 過期時間
				}
				token, _ := genToken(t, longSecret, cs, jwt.SigningMethodHS256)
				req.Header.Set("Authorization", "Bearer "+token)
			},
			wantHTTPStatus: http.StatusUnauthorized,
			wantJSON:       fmt.Sprintf(`{"error":"%s"}`, ErrTokenExpired),
			wantClaims:     nil,
		},
		{
			name:   "invalid token - modified token",
			config: AuthConfig{Secret: longSecret},
			setupHeader: func(req *http.Request) {
				cs := jwt.MapClaims{
					"iss": "test",
					"aud": "single",
					"exp": time.Now().Add(5 * time.Minute).Unix(), // 過期時間
				}
				token, _ := genToken(t, longSecret, cs, jwt.SigningMethodHS256)
				tokenBytes := []byte(token)
				modifiedToken := string(append(tokenBytes[:len(tokenBytes)-1], 'x')) // 故意修改最後一位
				req.Header.Set("Authorization", "Bearer "+modifiedToken)
			},
			wantHTTPStatus: http.StatusUnauthorized,
			wantJSON:       fmt.Sprintf(`{"error":"%s"}`, ErrSignatureInvalid),
			wantClaims:     nil,
		},
		{
			name:   "invalid token - algorithm mismatch",
			config: AuthConfig{Secret: longSecret},
			setupHeader: func(req *http.Request) {
				cs := jwt.MapClaims{
					"iss": "test",
					"aud": "single",
					"exp": time.Now().Add(5 * time.Minute).Unix(),
				}
				token := genWeiredToken(t, "none", cs, longSecret)
				req.Header.Set("Authorization", "Bearer "+token)
			},
			wantHTTPStatus: 401,
			wantJSON:       fmt.Sprintf(`{"error":"%s"}`, ErrUnsupportedAlgorithm),
			wantClaims:     nil,
		},
		{
			name:   "測試private claims",
			config: AuthConfig{Secret: longSecret},
			setupHeader: func(req *http.Request) {
				cs := jwt.MapClaims{
					"iss": "test",
					"aud": "single",
					"exp": time.Now().Add(5 * time.Minute).Unix(), // 過期時間
					// 私有聲明
					"aa": "aa_value",
					"bb": "bb_value",
					"cc": "cc_value",
				}
				token, _ := genToken(t, longSecret, cs, jwt.SigningMethodHS256)
				req.Header.Set("Authorization", "Bearer "+token)
			},
			wantHTTPStatus: 200,
			wantJSON:       "",
			wantClaims: &Claims[privateClaims]{
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "test",
					Audience:  jwt.ClaimStrings{"single"},
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
				},
				Private: privateClaims{
					AA: "aa_value",
					BB: "bb_value",
					CC: "cc_value",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			// 1. 建立 Gin engine
			engine := gin.New()

			// 2. 建立 Middleware

			mw, err := NewAuthMiddleware[privateClaims](tt.config)
			assert.NoError(t, err)
			engine.Use(mw.HandlerFunc())

			// 3. 建立一條測試路由
			var gotClaims *Claims[privateClaims]
			engine.GET("/test", func(ctx *gin.Context) {
				gotClaims = ctx.MustGet("claims").(*Claims[privateClaims])
				ctx.Status(http.StatusOK)
			})

			// 4. 準備 HTTP 請求
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			tt.setupHeader(req)

			// 5. Recorder
			respRec := httptest.NewRecorder()
			engine.ServeHTTP(respRec, req)

			// 6. 斷言
			assert.Equal(t, tt.wantHTTPStatus, respRec.Code)
			if tt.wantJSON != "" {
				assert.JSONEq(t, tt.wantJSON, respRec.Body.String())
			}
			if tt.wantClaims != nil {
				assert.Equal(t, tt.wantClaims, gotClaims)
			}
		})
	}
}
func genToken(t *testing.T, secret string, claims jwt.MapClaims, alg jwt.SigningMethod) (string, error) {
	t.Helper()
	token := jwt.NewWithClaims(alg, claims)
	return token.SignedString([]byte(secret))
}
func genWeiredToken(t *testing.T, weiredAlg string, claims jwt.MapClaims, secret string) string {
	t.Helper()

	// 1. header
	weiredHeader := fmt.Sprintf(`{"alg":"%s","typ":"JWT"}`, weiredAlg)
	headerSegment := base64.RawURLEncoding.EncodeToString(
		[]byte(weiredHeader),
	)

	// 2. payload
	payloadBytes, _ := json.Marshal(claims)
	payloadSegment := base64.RawURLEncoding.EncodeToString(payloadBytes)

	// 3. 用HS256 簽名騙過
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(headerSegment + "." + payloadSegment))
	signatureSegment := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	// 4. 組合成完整的 JWT
	return fmt.Sprintf("%s.%s.%s", headerSegment, payloadSegment, signatureSegment)

}
