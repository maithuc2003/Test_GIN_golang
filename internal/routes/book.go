package routes

import (
	"github.com/maithuc2003/Test_GIN_golang/internal/handler/book"
	RepInterface "github.com/maithuc2003/Test_GIN_golang/internal/interfaces/repositories"
	ServiceInterface "github.com/maithuc2003/Test_GIN_golang/internal/interfaces/service"
	Repo "github.com/maithuc2003/Test_GIN_golang/internal/repositories/book"
	ServiceImp "github.com/maithuc2003/Test_GIN_golang/internal/service/book"

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
