package book

import (
	"errors"
	"time"

	"github.com/maithuc2003/Test_GIN_golang/internal/interfaces/repositories"
	"github.com/maithuc2003/Test_GIN_golang/internal/models"
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
	books, err := s.bookRepo.GetAllBooks()
	if err != nil {
		return nil, err
	}

	if len(books) == 0 {
		return nil, errors.New("no books found")
	}

	return books, nil
}
