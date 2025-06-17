package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	RegisterBookRoutes(r, db)
	RegisterUserRoutes(r, db)

	return r
}
