package book

import (
	"net/http"

	"github.com/maithuc2003/Test_GIN_golang/internal/interfaces/service"
	"github.com/maithuc2003/Test_GIN_golang/internal/models"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	bookService service.BookServiceInterface
}

func NewBookHandler(bookService service.BookServiceInterface) *BookHandler {
	return &BookHandler{bookService: bookService}
}

// POST/ books
func (h *BookHandler) CreateBookHandler(c *gin.Context) {

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.bookService.CreateBook(&book)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
	}
	c.JSON(http.StatusCreated, book)
}

// Get/books
func (h *BookHandler) GetAllBooksHandler(c *gin.Context) {
	books, err := h.bookService.GetAllBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}
	// c.JSON(http.StatusOK, books)
	c.IndentedJSON(http.StatusOK, books)

}
