package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maithuc2003/Test_GIN_golang/internal/handler/book"
	RepInterface "github.com/maithuc2003/Test_GIN_golang/internal/interfaces/repositories"
	ServiceInterface "github.com/maithuc2003/Test_GIN_golang/internal/interfaces/service"
	"github.com/maithuc2003/Test_GIN_golang/internal/middleware"
	Repo "github.com/maithuc2003/Test_GIN_golang/internal/repositories/book"
	ServiceImp "github.com/maithuc2003/Test_GIN_golang/internal/service/book"
	"gorm.io/gorm"
)

func RegisterBookRoutes(r *gin.Engine, db *gorm.DB) {
	var bookRepo RepInterface.BookRepository = Repo.NewRepository(db)
	var bookService ServiceInterface.BookServiceInterface = ServiceImp.NewBookService(bookRepo)
	bookHandler := book.NewBookHandler(bookService)

	// Public routes
	bookRoutes := r.Group("/books")
	{
		bookRoutes.GET("", bookHandler.GetAllBooksHandler)
		bookRoutes.GET("/:id", bookHandler.GetByBookID)
	}

	// Protected routes with Auth + RBAC
	auth := r.Group("/books", middleware.AuthMiddleware())
	{
		auth.POST("/add", middleware.RBACMiddleware("book/create"), bookHandler.CreateBookHandler)
		auth.PUT("/:id", middleware.RBACMiddleware("book/update"), bookHandler.UpdateById)
		auth.DELETE("/:id", middleware.RBACMiddleware("book/delete"), bookHandler.DeleteById)
	}
}
