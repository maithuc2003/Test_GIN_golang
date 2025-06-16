package mocks

import (
	"github.com/maithuc2003/Test_GIN_golang/internal/models"
	"github.com/stretchr/testify/mock"
)

// 👉 MockBookService implements service.BookService
type MockBookService struct {
	mock.Mock
}

func (m *MockBookService) CreateBook(book *models.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookService) GetAllBooks() ([]models.Book, error) {
	args := m.Called()
	return args.Get(0).([]models.Book), args.Error(1)
}
