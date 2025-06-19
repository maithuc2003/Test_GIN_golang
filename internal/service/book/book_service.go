package book

import (
	"errors"
	"strings"
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
func (s *BookService) GetByBookID(id int) (*models.Book, error) {
	if id <= 0 {
		return nil, errors.New("invalid book ID")
	}
	return s.bookRepo.GetByBookID(id)
}

func (s *BookService) DeleteById(id int) (*models.Book, error) {
	if id <= 0 {
		return nil, errors.New("invalid book ID")
	}
	return s.bookRepo.DeleteById(id)
}

func (s *BookService) UpdateById(book *models.Book) (*models.Book, error) {
	if book == nil {
		return nil, errors.New("book is nil")
	}
	if book.ID <= 0 {
		return nil, errors.New("invalid book ID")
	}
	if strings.TrimSpace(book.Title) == "" {
		return nil, errors.New("book title is required")
	}
	if book.AuthorID <= 0 {
		return nil, errors.New("book author ID is required")
	}
	if book.Stock < 0 {
		return nil, errors.New("book quantity cannot be negative")
	}

	return s.bookRepo.UpdateById(book)
}
