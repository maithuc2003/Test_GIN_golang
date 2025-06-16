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
	mocks "github.com/maithuc2003/Test_GIN_golang/internal/handler/book/mocks"
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
			mockReturnErr:  nil, // won't reach service layer
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			inputBody: models.Book{
				Title:    "Book 2",
				Stock:    5,
				AuthorID: 2,
			},
			mockReturnErr:  errors.New("DB error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt // capture
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(mocks.MockBookService)
			h := book.NewBookHandler(mockService)

			var reqBody []byte
			switch body := tt.inputBody.(type) {
			case string:
				reqBody = []byte(body)
			default:
				reqBody, _ = json.Marshal(body)
			}

			if _, ok := tt.inputBody.(models.Book); ok {
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
			name: "success fetch",
			mockBooks: []models.Book{
				{ID: 1, Title: "Book 1", Stock: 5, AuthorID: 1},
				{ID: 2, Title: "Book 2", Stock: 10, AuthorID: 2},
			},
			mockReturnErr:  nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "error from service",
			mockBooks:      nil,
			mockReturnErr:  errors.New("fetch error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		tt := tt
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

			if rec.Code == http.StatusOK {
				var got []models.Book
				err := json.Unmarshal(rec.Body.Bytes(), &got)
				require.NoError(t, err)
				require.Equal(t, tt.mockBooks, got)
			} else {
				var got map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &got)
				require.NoError(t, err)
				require.Contains(t, got["error"], "Failed to fetch books")
			}

			mockService.AssertExpectations(t)
		})
	}
}
