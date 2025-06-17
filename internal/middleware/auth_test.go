package middleware_test

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"

	"github.com/maithuc2003/Test_GIN_golang/internal/middleware"
	jwtutil "github.com/maithuc2003/Test_GIN_golang/pkg/jwt"
)

// Tạo token hợp lệ
func generateToken(userID int, username, role string, expired bool) string {
	expTime := time.Now().Add(time.Hour)
	if expired {
		expTime = time.Now().Add(-time.Hour)
	}

	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"aud":      role,
		"exp":      expTime.Unix(),
		"iat":      time.Now().Unix(),
		"iss":      "maithuc",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString(jwtutil.JwtSecret())
	return signedToken
}

func setupRouterWithMiddleware() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.AuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
	})
	return r
}
func generateRSAToken(t *testing.T) string {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	claims := jwt.MapClaims{
		"user_id":  1,
		"username": "attacker",
		"aud":      "admin",
		"exp":      time.Now().Add(time.Hour).Unix(),
		"iat":      time.Now().Unix(),
		"iss":      "fake",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	signedToken, err := token.SignedString(privateKey)
	require.NoError(t, err)

	return signedToken
}

func TestAuthMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Missing token",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Missing or invalid token",
		},
		{
			name:           "Invalid format token",
			authHeader:     "InvalidToken",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Missing or invalid token",
		},
		{
			name:           "Invalid signature token",
			authHeader:     "Bearer invalid.token.value",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Token is invalid or expired",
		},
		{
			name:           "Valid token but not admin",
			authHeader:     "Bearer " + generateToken(1, "user", "user", false),
			expectedStatus: http.StatusForbidden,
			expectedBody:   "Access denied",
		},
		{
			name:           "Valid admin token",
			authHeader:     "Bearer " + generateToken(1, "admin", "admin", false),
			expectedStatus: http.StatusOK,
			expectedBody:   "Success",
		},
		{
			name:           "Expired token",
			authHeader:     "Bearer " + generateToken(1, "admin", "admin", true),
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Token is invalid or expired",
		},
		{
			name:           "Invalid signing method",
			authHeader:     "Bearer " + generateRSAToken(t), // Dùng RSA
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Token is invalid or expired",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			router := setupRouterWithMiddleware()
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if tc.authHeader != "" {
				req.Header.Set("auth", tc.authHeader)
			}
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			require.Equal(t, tc.expectedStatus, w.Code)
			require.Contains(t, w.Body.String(), tc.expectedBody)
		})
	}
}
