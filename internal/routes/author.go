package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maithuc2003/Test_GIN_golang/internal/handler/author"
	RepInterface "github.com/maithuc2003/Test_GIN_golang/internal/interfaces/repositories"
	ServiceInterface "github.com/maithuc2003/Test_GIN_golang/internal/interfaces/service"
	"github.com/maithuc2003/Test_GIN_golang/internal/middleware"
	Repo "github.com/maithuc2003/Test_GIN_golang/internal/repositories/author"
	ServiceImp "github.com/maithuc2003/Test_GIN_golang/internal/service/author"
	"gorm.io/gorm"
)

func RegisterAuthorRoutes(r *gin.Engine, db *gorm.DB) {
	var authorRepo RepInterface.AuthorRepositoriesInterface = Repo.NewAuthorRepo(db)
	var authorService ServiceInterface.AuthorServiceInterface = ServiceImp.NewAuthorService(authorRepo)
	authorHandler := author.NewAuthorHandler(authorService)

	// Public routes
	r.GET("/authors", authorHandler.GetAllAuthors)
	r.GET("/authors/:id", authorHandler.GetByAuthorID)

	// Authenticated (admin) routes
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.POST("/authors", authorHandler.CreateAuthor)
	auth.PUT("/authors/:id", authorHandler.UpdateById)
	auth.DELETE("/authors/:id", authorHandler.DeleteById)
}
