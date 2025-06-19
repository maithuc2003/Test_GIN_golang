package book

import (
	"net/http"
	"strconv"
	"time"

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

// GET /books/:id
func (h *BookHandler) GetByBookID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	book, err := h.bookService.GetByBookID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

// DELETE /books/:id
func (h *BookHandler) DeleteById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	book, err := h.bookService.DeleteById(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Delete failed or book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}

// PUT /books/:id
func (h *BookHandler) UpdateById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var updateBook models.Book
	if err := c.ShouldBindJSON(&updateBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updateBook.ID = uint(id)
	updateBook.UpdatedAt = time.Now()

	book, err := h.bookService.UpdateById(&updateBook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}
	c.JSON(http.StatusOK, book)
}
