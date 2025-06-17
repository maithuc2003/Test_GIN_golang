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

	// r.POST("/book/add", bookHandler.CreateBookHandler)
	r.GET("/books", bookHandler.GetAllBooksHandler)

	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.POST("/book/add", bookHandler.CreateBookHandler) // Chỉ admin được gọi
	

}
