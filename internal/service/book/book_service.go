package book

import (
	"github.com/maithuc2003/GIN_golang_framework/internal/interfaces/repositories"
	"github.com/maithuc2003/GIN_golang_framework/internal/models"
	"errors"
	"time"
)

type BookService struct {
	bookRepo repositories.BookRepository
}

func NewBookService(bookRepo repositories.BookRepository) *BookService {
	return &BookService{bookRepo: bookRepo}

}

func (s *BookService) CreateBook(book *models.Book) error {
	if book.Title == "" || book.AuthorID == 0 {
		return errors.New("invalid book data: title and author_id required")
	}
	book.CreatedAt = time.Now()
	book.UpdatedAt = time.Now()
	return s.bookRepo.CreateBook(book)
}

func (s *BookService) GetAllBooks() ([]models.Book, error) {
	return s.bookRepo.GetAllBooks()
}
