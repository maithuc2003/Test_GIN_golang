package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	jwtutil "github.com/maithuc2003/Test_GIN_golang/pkg/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("auth")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid token"})
			c.Abort()
			return
		}

		// Cắt "Bearer " để lấy token thật sự
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse token và xác minh chữ ký
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return jwtutil.JwtSecret(), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is invalid or expired"})
			c.Abort()
			return
		}

		// Lấy thông tin từ token và lưu vào context
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			if userID, ok := claims["user_id"].(float64); ok {
				c.Set("user_id", int(userID))
			}
			if username, ok := claims["username"].(string); ok {
				c.Set("username", username)
			}
			if role, ok := claims["aud"].(string); ok {
				c.Set("role", role)
				if role != "admin" {
					c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}
