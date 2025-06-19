package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/maithuc2003/Test_GIN_golang/internal/database"
)

func RBACMiddleware(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id") // phải trùng key AuthMiddleware gán
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
			c.Abort()
			return
		}

		if !database.HasAccess(userID.(int), permission) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Permission denied"})
			c.Abort()
			return
		}
		c.Next()
	}
}
