package book_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mocks "github.com/maithuc2003/Test_GIN_golang/internal/mocks/repositories"
	"github.com/maithuc2003/Test_GIN_golang/internal/models"
	"github.com/maithuc2003/Test_GIN_golang/internal/service/book"
)

func TestCreateBook(t *testing.T) {
	mockRepo := new(mocks.MockBookRepo)
	service := book.NewBookService(mockRepo)

	tests := []struct {
		name        string
		input       *models.Book
		setupMock   func()
		expectError bool
	}{
		{
			name: "missing title and author_id",
			input: &models.Book{
				Title:    "",
				AuthorID: 0,
			},
			setupMock:   func() {},
			expectError: true,
		},
		{
			name: "valid book",
			input: &models.Book{
				Title:    "Go Book",
				AuthorID: 1,
			},
			setupMock: func() {
				mockRepo.On("CreateBook", mock.AnythingOfType("*models.Book")).Return(nil).Once()
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			err := service.CreateBook(tt.input)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetAllBooks(t *testing.T) {
	mockRepo := new(mocks.MockBookRepo)
	service := book.NewBookService(mockRepo)

	tests := []struct {
		name        string
		mockReturn  []models.Book
		mockError   error
		expectError bool
	}{
		{
			name:        "empty result",
			mockReturn:  []models.Book{},
			mockError:   nil,
			expectError: true,
		},
		{
			name: "books found",
			mockReturn: []models.Book{
				{ID: 1, Title: "Go", AuthorID: 1},
			},
			mockError:   nil,
			expectError: false,
		},
		{
			name:        "repo error",
			mockReturn:  nil,
			mockError:   errors.New("database error"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.On("GetAllBooks").Return(tt.mockReturn, tt.mockError).Once()
			result, err := service.GetAllBooks()
			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.mockReturn, result)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetByBookID(t *testing.T) {
	mockRepo := new(mocks.MockBookRepo)
	service := book.NewBookService(mockRepo)

	tests := []struct {
		name        string
		inputID     int
		mockReturn  *models.Book
		mockError   error
		expectError bool
	}{
		{
			name:        "invalid ID",
			inputID:     0,
			expectError: true,
		},
		{
			name:        "valid ID",
			inputID:     1,
			mockReturn:  &models.Book{ID: 1, Title: "Go"},
			mockError:   nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.inputID > 0 {
				mockRepo.On("GetByBookID", tt.inputID).Return(tt.mockReturn, tt.mockError).Once()
			}
			result, err := service.GetByBookID(tt.inputID)
			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.mockReturn, result)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteById(t *testing.T) {
	mockRepo := new(mocks.MockBookRepo)
	service := book.NewBookService(mockRepo)

	tests := []struct {
		name        string
		inputID     int
		mockReturn  *models.Book
		mockError   error
		expectError bool
	}{
		{
			name:        "invalid ID",
			inputID:     -1,
			expectError: true,
		},
		{
			name:        "valid delete",
			inputID:     1,
			mockReturn:  &models.Book{ID: 1},
			mockError:   nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.inputID > 0 {
				mockRepo.On("DeleteById", tt.inputID).Return(tt.mockReturn, tt.mockError).Once()
			}
			result, err := service.DeleteById(tt.inputID)
			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.mockReturn, result)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateById(t *testing.T) {
	mockRepo := new(mocks.MockBookRepo)
	service := book.NewBookService(mockRepo)

	tests := []struct {
		name        string
		input       *models.Book
		mockReturn  *models.Book
		mockError   error
		expectError bool
	}{
		{
			name:        "nil book",
			input:       nil,
			expectError: true,
		},
		{
			name:        "invalid ID",
			input:       &models.Book{ID: 0, Title: "Go", AuthorID: 1},
			expectError: true,
		},
		{
			name:        "empty title",
			input:       &models.Book{ID: 1, Title: "", AuthorID: 1},
			expectError: true,
		},
		{
			name:        "missing author_id",
			input:       &models.Book{ID: 1, Title: "Go", AuthorID: 0},
			expectError: true,
		},
		{
			name:        "negative stock",
			input:       &models.Book{ID: 1, Title: "Go", AuthorID: 1, Stock: -10},
			expectError: true,
		},
		{
			name: "valid update",
			input: &models.Book{
				ID:       1,
				Title:    "Go",
				AuthorID: 1,
				Stock:    5,
			},
			mockReturn:  &models.Book{ID: 1, Title: "Go"},
			mockError:   nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input != nil && tt.input.ID > 0 && tt.input.AuthorID > 0 && tt.input.Stock >= 0 && tt.input.Title != "" {
				mockRepo.On("UpdateById", tt.input).Return(tt.mockReturn, tt.mockError).Once()
			}
			result, err := service.UpdateById(tt.input)
			if tt.expectError {
				require.Error(t, err)
				require.Nil(t, result)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.mockReturn, result)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
