package mocks

import (
	"github.com/maithuc2003/GIN_golang_framework/internal/models"

	"github.com/stretchr/testify/mock"
)

type MockBookRepo struct {
	mock.Mock
}

func (m *MockBookRepo) CreateBook(book *models.Book) error {
	args := m.Called(book)
	return args.Error(0)
}

func (m *MockBookRepo) GetAllBooks() ([]models.Book, error) {
	args := m.Called()
	return args.Get(0).([]models.Book), args.Error(1)
}
