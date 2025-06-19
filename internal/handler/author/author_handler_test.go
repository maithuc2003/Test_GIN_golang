package author_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/maithuc2003/Test_GIN_golang/internal/handler/author"
	mockService "github.com/maithuc2003/Test_GIN_golang/internal/mocks/service"
	"github.com/maithuc2003/Test_GIN_golang/internal/models"
)

func TestGetAllAuthors(t *testing.T) {
	type testCase struct {
		name       string
		mockData   []*models.Author
		mockErr    error
		wantStatus int
	}

	tests := []testCase{
		{
			name: "Success",
			mockData: []*models.Author{
				{ID: 1, Name: "Author A"},
			},
			mockErr:    nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "InternalError",
			mockData:   nil,
			mockErr:    errors.New("failed"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := new(mockService.MockAuthorService)
			mockSvc.On("GetAllAuthors").Return(tc.mockData, tc.mockErr)

			r := gin.Default()
			handler := author.NewAuthorHandler(mockSvc)
			r.GET("/authors", handler.GetAllAuthors)

			req, _ := http.NewRequest("GET", "/authors", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tc.wantStatus, w.Code)
		})
	}
}
func TestCreateAuthor(t *testing.T) {
	type testCase struct {
		name       string
		input      *models.Author
		rawBody    string // dùng cho test Invalid JSON
		mockErr    error
		wantStatus int
	}

	tests := []testCase{
		{
			name: "Success",
			input: &models.Author{
				Name: "Author A",
			},
			mockErr:    nil,
			wantStatus: http.StatusCreated,
		},
		{
			name:       "Invalid JSON",
			rawBody:    "invalid json",
			mockErr:    nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "Create Failed",
			input: &models.Author{
				Name: "Error Author",
			},
			mockErr:    errors.New("failed to create"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := new(mockService.MockAuthorService)
			handler := author.NewAuthorHandler(mockSvc)

			r := gin.Default()
			r.POST("/authors", handler.CreateAuthor)

			var req *http.Request
			if tc.rawBody != "" {
				req, _ = http.NewRequest("POST", "/authors", bytes.NewBufferString(tc.rawBody))
			} else {
				tc.input.CreatedAt = time.Now()
				body, _ := json.Marshal(tc.input)
				req, _ = http.NewRequest("POST", "/authors", bytes.NewBuffer(body))
				mockSvc.On("CreateAuthor", mock.Anything).Return(tc.mockErr)
			}

			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			require.Equal(t, tc.wantStatus, w.Code)
		})
	}
}

func TestGetByAuthorID(t *testing.T) {
	type testCase struct {
		name       string
		param      string
		mockData   *models.Author
		mockErr    error
		wantStatus int
	}

	tests := []testCase{
		{
			name:       "Success",
			param:      "1",
			mockData:   &models.Author{ID: 1, Name: "Author"},
			mockErr:    nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid ID",
			param:      "abc",
			mockData:   nil,
			mockErr:    nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Error from service",
			param:      "2",
			mockData:   nil,
			mockErr:    errors.New("not found"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := new(mockService.MockAuthorService)
			if tc.mockData != nil || tc.mockErr != nil {
				mockSvc.On("GetByAuthorID", mock.Anything).Return(tc.mockData, tc.mockErr)
			}

			r := gin.Default()
			handler := author.NewAuthorHandler(mockSvc)
			r.GET("/authors/:id", handler.GetByAuthorID)

			req, _ := http.NewRequest("GET", "/authors/"+tc.param, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tc.wantStatus, w.Code)
		})
	}
}
func TestDeleteById(t *testing.T) {
	type testCase struct {
		name       string
		param      string
		mockData   *models.Author
		mockErr    error
		wantStatus int
	}

	tests := []testCase{
		{
			name:       "Success",
			param:      "1",
			mockData:   &models.Author{ID: 1, Name: "Deleted"},
			mockErr:    nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid ID",
			param:      "abc",
			mockData:   nil,
			mockErr:    nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Service Error",
			param:      "2",
			mockData:   nil,
			mockErr:    errors.New("delete failed"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := new(mockService.MockAuthorService)
			if tc.wantStatus != http.StatusBadRequest {
				mockSvc.On("DeleteById", mock.Anything).Return(tc.mockData, tc.mockErr)
			}

			r := gin.Default()
			handler := author.NewAuthorHandler(mockSvc)
			r.DELETE("/authors/:id", handler.DeleteById)

			req, _ := http.NewRequest("DELETE", "/authors/"+tc.param, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			require.Equal(t, tc.wantStatus, w.Code)
		})
	}
}
func TestUpdateById(t *testing.T) {
	type testCase struct {
		name       string
		param      string
		input      *models.Author
		rawBody    string
		mockResult *models.Author
		mockErr    error
		wantStatus int
	}

	tests := []testCase{
		{
			name:  "Success",
			param: "1",
			input: &models.Author{Name: "Updated"},
			mockResult: &models.Author{
				ID:   1,
				Name: "Updated",
			},
			mockErr:    nil,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid ID Param",
			param:      "abc",
			input:      nil,
			mockResult: nil,
			mockErr:    nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Invalid JSON",
			param:      "1",
			rawBody:    "not-json",
			mockResult: nil,
			mockErr:    nil,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Service Error",
			param:      "2",
			input:      &models.Author{Name: "Service Error"},
			mockResult: nil,
			mockErr:    errors.New("update failed"),
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockSvc := new(mockService.MockAuthorService)
			handler := author.NewAuthorHandler(mockSvc)

			r := gin.Default()
			r.PUT("/authors/:id", handler.UpdateById)

			var req *http.Request
			if tc.rawBody != "" {
				req, _ = http.NewRequest("PUT", "/authors/"+tc.param, bytes.NewBufferString(tc.rawBody))
			} else if tc.input != nil {
				tc.input.UpdatedAt = time.Now()
				body, _ := json.Marshal(tc.input)
				req, _ = http.NewRequest("PUT", "/authors/"+tc.param, bytes.NewBuffer(body))

				if tc.wantStatus != http.StatusBadRequest {
					mockSvc.On("UpdateById", mock.Anything).Return(tc.mockResult, tc.mockErr)
				}
			} else {
				req, _ = http.NewRequest("PUT", "/authors/"+tc.param, nil)
			}

			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			require.Equal(t, tc.wantStatus, w.Code)
		})
	}
}
