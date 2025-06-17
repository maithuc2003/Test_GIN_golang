package jwtutil

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = func() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// fallback default for local dev/test
		secret = "default_secret"
	}
	return []byte(secret)
}()

func JwtSecret() []byte {
	return jwtSecret
}
func GenerateJWT(userID uint, username string, role string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"sub":      username,                              // Subject
		"iss":      "maithuc",                             // Issuer
		"aud":      role,                                  // Audience (ví dụ: "admin" hoặc "user") role permission
		"exp":      time.Now().Add(72 * time.Hour).Unix(), // Expiration
		"iat":      time.Now().Unix(),                     // Issued at book action order action (list , add )
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
