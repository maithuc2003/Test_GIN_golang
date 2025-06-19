package mocks

import (
	"github.com/maithuc2003/Test_GIN_golang/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockAuthorRepository struct {
	mock.Mock
}

// GetAllAuthors mocks retrieving all authors
func (m *MockAuthorRepository) GetAllAuthors() ([]*models.Author, error) {
	args := m.Called()
	if authors, ok := args.Get(0).([]*models.Author); ok {
		return authors, args.Error(1)
	}
	return nil, args.Error(1)
}

// CreateAuthor mocks creating a new author
func (m *MockAuthorRepository) CreateAuthor(author *models.Author) error {
	args := m.Called(author)
	return args.Error(0)
}

// GetByAuthorID mocks retrieving an author by ID
func (m *MockAuthorRepository) GetByAuthorID(id int) (*models.Author, error) {
	args := m.Called(id)
	if author, ok := args.Get(0).(*models.Author); ok {
		return author, args.Error(1)
	}
	return nil, args.Error(1)
}

// UpdateById mocks updating an author
func (m *MockAuthorRepository) UpdateById(author *models.Author) (*models.Author, error) {
	args := m.Called(author)
	if updated, ok := args.Get(0).(*models.Author); ok {
		return updated, args.Error(1)
	}
	return nil, args.Error(1)
}

// DeleteById mocks deleting an author by ID
func (m *MockAuthorRepository) DeleteById(id int) (*models.Author, error) {
	args := m.Called(id)
	if deleted, ok := args.Get(0).(*models.Author); ok {
		return deleted, args.Error(1)
	}
	return nil, args.Error(1)
}
