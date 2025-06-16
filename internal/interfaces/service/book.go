package service

import "github.com/maithuc2003/Test_GIN_golang/internal/models"

type BookServiceInterface interface {
	CreateBook(book *models.Book) error
	GetAllBooks() ([]models.Book, error)
}
