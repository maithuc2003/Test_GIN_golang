package book_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maithuc2003/Test_GIN_golang/internal/handler/book"
	mocks "github.com/maithuc2003/Test_GIN_golang/internal/mocks/service"
	"github.com/maithuc2003/Test_GIN_golang/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateBookHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		inputBody      interface{}
		mockReturnErr  error
		expectedStatus int
	}{
		{
			name: "valid input",
			inputBody: models.Book{
				Title:    "Book 1",
				Stock:    10,
				AuthorID: 1,
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid JSON",
			inputBody:      `{"title": "Book 1", "stock": "wrong_type"}`,
			mockReturnErr:  nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			inputBody: models.Book{
				Title:    "Book 2",
				Stock:    5,
				AuthorID: 2,
			},
			mockReturnErr:  errors.New("failed to create"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.MockBookService)
			h := book.NewBookHandler(mockService)

			var reqBody []byte
			switch v := tt.inputBody.(type) {
			case string:
				reqBody = []byte(v)
			default:
				reqBody, _ = json.Marshal(v)
				mockService.On("CreateBook", mock.AnythingOfType("*models.Book")).Return(tt.mockReturnErr)
			}

			req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			r := gin.Default()
			r.POST("/books", h.CreateBookHandler)
			r.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			mockService.AssertExpectations(t)
		})
	}
}

func TestGetAllBooksHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockBooks      []models.Book
		mockReturnErr  error
		expectedStatus int
	}{
		{
			name: "success",
			mockBooks: []models.Book{
				{ID: 1, Title: "Book A", AuthorID: 1, Stock: 5},
				{ID: 2, Title: "Book B", AuthorID: 2, Stock: 10},
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "service error",
			mockBooks:      nil,
			mockReturnErr:  errors.New("fetch error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.MockBookService)
			h := book.NewBookHandler(mockService)

			mockService.On("GetAllBooks").Return(tt.mockBooks, tt.mockReturnErr)

			req := httptest.NewRequest(http.MethodGet, "/books", nil)
			rec := httptest.NewRecorder()

			r := gin.Default()
			r.GET("/books", h.GetAllBooksHandler)
			r.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			mockService.AssertExpectations(t)

			if tt.expectedStatus == http.StatusOK {
				var got []models.Book
				err := json.Unmarshal(rec.Body.Bytes(), &got)
				require.NoError(t, err)
				require.Equal(t, tt.mockBooks, got)
			}
		})
	}
}

func TestGetByBookIDHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		paramID        string
		mockBook       *models.Book
		mockReturnErr  error
		expectedStatus int
	}{
		{
			name:           "valid ID",
			paramID:        "1",
			mockBook:       &models.Book{ID: 1, Title: "Book A", AuthorID: 1, Stock: 5},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid ID format",
			paramID:        "abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "book not found",
			paramID:        "99",
			mockReturnErr:  errors.New("not found"),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.MockBookService)
			h := book.NewBookHandler(mockService)

			if tt.expectedStatus != http.StatusBadRequest {
				mockService.On("GetByBookID", mock.AnythingOfType("int")).Return(tt.mockBook, tt.mockReturnErr)
			}

			req := httptest.NewRequest(http.MethodGet, "/books/"+tt.paramID, nil)
			rec := httptest.NewRecorder()

			r := gin.Default()
			r.GET("/books/:id", h.GetByBookID)
			r.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			mockService.AssertExpectations(t)
		})
	}
}
func TestUpdateBookHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		paramID        string
		inputBody      interface{}
		mockReturnBook *models.Book
		mockReturnErr  error
		expectedStatus int
	}{
		{
			name:    "valid update",
			paramID: "1",
			inputBody: models.Book{
				Title:    "Updated Title",
				AuthorID: 1,
				Stock:    15,
			},
			mockReturnBook: &models.Book{ID: 1, Title: "Updated Title", AuthorID: 1, Stock: 15},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid ID param",
			paramID:        "abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid JSON body",
			paramID:        "1",
			inputBody:      `{"title": 123}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "service error",
			paramID: "2",
			inputBody: models.Book{
				Title:    "Error Book",
				AuthorID: 2,
				Stock:    10,
			},
			mockReturnErr:  errors.New("update failed"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.MockBookService)
			h := book.NewBookHandler(mockService)

			var reqBody []byte
			switch v := tt.inputBody.(type) {
			case string:
				reqBody = []byte(v)
			case nil:
				reqBody = nil
			default:
				reqBody, _ = json.Marshal(v)
				mockService.On("UpdateById", mock.AnythingOfType("*models.Book")).Return(tt.mockReturnBook, tt.mockReturnErr)
			}

			req := httptest.NewRequest(http.MethodPut, "/books/"+tt.paramID, bytes.NewBuffer(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			r := gin.Default()
			r.PUT("/books/:id", h.UpdateById)
			r.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			mockService.AssertExpectations(t)
		})
	}
}
func TestDeleteBookHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		paramID        string
		mockReturnBook *models.Book
		mockReturnErr  error
		expectedStatus int
	}{
		{
			name:           "valid delete",
			paramID:        "1",
			mockReturnBook: &models.Book{ID: 1, Title: "To Delete"},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid ID param",
			paramID:        "abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "book not found",
			paramID:        "99",
			mockReturnErr:  errors.New("not found"),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.MockBookService)
			h := book.NewBookHandler(mockService)

			if tt.expectedStatus != http.StatusBadRequest {
				mockService.On("DeleteById", mock.AnythingOfType("int")).Return(tt.mockReturnBook, tt.mockReturnErr)
			}

			req := httptest.NewRequest(http.MethodDelete, "/books/"+tt.paramID, nil)
			rec := httptest.NewRecorder()

			r := gin.Default()
			r.DELETE("/books/:id", h.DeleteById)
			r.ServeHTTP(rec, req)

			require.Equal(t, tt.expectedStatus, rec.Code)
			mockService.AssertExpectations(t)
		})
	}
}
