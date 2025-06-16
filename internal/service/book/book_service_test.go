package book_test

import (
	"errors"
	"testing"
	"time"

	"github.com/maithuc2003/Test_GIN_golang/internal/interfaces/service"
	"github.com/maithuc2003/Test_GIN_golang/internal/models"
	bookImpl "github.com/maithuc2003/Test_GIN_golang/internal/service/book"
	"github.com/maithuc2003/Test_GIN_golang/internal/service/book/mocks"

	"github.com/stretchr/testify/assert"
)

func TestCreateBook(t *testing.T) {
	tests := []struct {
		name        string
		input       *models.Book
		mockReturn  error
		expectError bool
	}{
		{
			name: "valid book",
			input: &models.Book{
				Title:    "Clean Code",
				AuthorID: 1,
			},
			mockReturn:  nil,
			expectError: false,
		},
		{
			name: "missing title",
			input: &models.Book{
				Title:    "",
				AuthorID: 1,
			},
			expectError: true,
		},
		{
			name: "missing author",
			input: &models.Book{
				Title:    "Some Book",
				AuthorID: 0,
			},
			expectError: true,
		},
		{
			name: "repo returns error",
			input: &models.Book{
				Title:    "Failing Save",
				AuthorID: 2,
			},
			mockReturn:  errors.New("DB error"),
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockBookRepo)

			// Only call repo if input is valid
			if tc.input.Title != "" && tc.input.AuthorID != 0 {
				mockRepo.On("CreateBook", tc.input).Return(tc.mockReturn)
			}

			var svc service.BookServiceInterface = bookImpl.NewBookService(mockRepo)
			err := svc.CreateBook(tc.input)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.WithinDuration(t, time.Now(), tc.input.CreatedAt, time.Second)
				assert.WithinDuration(t, time.Now(), tc.input.UpdatedAt, time.Second)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetAllBooks(t *testing.T) {
	tests := []struct {
		name         string
		mockBooks    []models.Book
		mockError    error
		expectError  bool
		expectedLen  int
		expectedMsg  string
	}{
		{
			name: "success - return books",
			mockBooks: []models.Book{
				{ID: 1, Title: "Clean Code", AuthorID: 1},
				{ID: 2, Title: "Go Programming", AuthorID: 2},
			},
			mockError:   nil,
			expectError: false,
			expectedLen: 2,
		},
		{
			name:         "repo error",
			mockBooks:    nil,
			mockError:    errors.New("db error"),
			expectError:  true,
			expectedMsg:  "db error",
		},
		{
			name:         "no books found",
			mockBooks:    []models.Book{},
			mockError:    nil,
			expectError:  true,
			expectedMsg:  "no books found",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(mocks.MockBookRepo)
			mockRepo.On("GetAllBooks").Return(tc.mockBooks, tc.mockError)

			var svc service.BookServiceInterface = bookImpl.NewBookService(mockRepo)
			books, err := svc.GetAllBooks()

			if tc.expectError {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedMsg)
			} else {
				assert.NoError(t, err)
				assert.Len(t, books, tc.expectedLen)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
