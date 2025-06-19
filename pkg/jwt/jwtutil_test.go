package jwtutil_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	jwtutil "github.com/maithuc2003/Test_GIN_golang/pkg/jwt"
	"github.com/stretchr/testify/require"
)

func TestGenerateJWT_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		userID   uint
		username string
	}{
		{"Valid user 1", 1, "john"},
		{"Valid user 99", 99, "maithuc2003"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			tokenStr, err := jwtutil.GenerateJWT(tt.userID, tt.username)

			// Assert
			require.NoError(t, err)
			require.NotEmpty(t, tokenStr)

			// Parse token
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					t.Fatalf("Unexpected signing method: %v", token.Header["alg"])
				}
				return jwtutil.JwtSecret(), nil
			})
			require.NoError(t, err)
			require.True(t, token.Valid)

			// Check claims
			claims, ok := token.Claims.(jwt.MapClaims)
			require.True(t, ok)

			require.EqualValues(t, tt.userID, int(claims["user_id"].(float64)))
			require.Equal(t, tt.username, claims["username"])

			// Check expiration time is roughly 72 hours from now
			require.WithinDuration(t,
				time.Now().Add(72*time.Hour),
				time.Unix(int64(claims["exp"].(float64)), 0),
				time.Minute,
			)
		})
	}
}
