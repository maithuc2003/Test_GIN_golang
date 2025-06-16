package routes

import (
	"github.com/maithuc2003/GIN_golang_framework/internal/handler/book"
	RepInterface "github.com/maithuc2003/GIN_golang_framework/internal/interfaces/repositories"
	ServiceInterface "github.com/maithuc2003/GIN_golang_framework/internal/interfaces/service"
	Repo "github.com/maithuc2003/GIN_golang_framework/internal/repositories"
	ServiceImp "github.com/maithuc2003/GIN_golang_framework/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	// 1. Khởi tạo Gin router
	r := gin.Default()
	// 2. Khởi tạo Repository
	var bookRepo RepInterface.BookRepository = Repo.NewRepository(db)

	// 3. Khảo tạo service, inject Repo
	var bookService ServiceInterface.BookServiceInterface = ServiceImp.NewBookService(bookRepo)

	//4. Khởi tạo Handler, inject Service
	BookHandler := book.NewBookHandler(bookService)

	r.POST("/book/add", BookHandler.CreateBookHandler)
	r.GET("/books", BookHandler.GetAllBooksHandler)
	return r
}
