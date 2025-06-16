package repositories

import "github.com/maithuc2003/Test_GIN_golang/internal/models"

type BookRepository interface {
	CreateBook(book *models.Book) error
	GetAllBooks() ([]models.Book, error)
}
