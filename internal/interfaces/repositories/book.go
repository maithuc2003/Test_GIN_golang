package repositories

import "github.com/maithuc2003/GIN_golang_framework/internal/models"

type BookRepository interface {
	CreateBook(book *models.Book) error
	GetAllBooks() ([]models.Book, error)
}
