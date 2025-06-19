package order_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/maithuc2003/Test_GIN_golang/internal/handler/order"
	mockService "github.com/maithuc2003/Test_GIN_golang/internal/mocks/service"
	"github.com/maithuc2003/Test_GIN_golang/internal/models"
)

func TestCreateOrder(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		input          interface{}
		mockErr        error
		expectedStatus int
	}{
		{
			name: "valid",
			input: models.Order{
				BookID:   1,
				UserID:   1,
				Quantity: 2,
				Status:   "pending",
			},
			mockErr:        nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid JSON",
			input:          `{"invalid":::}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			input: models.Order{
				BookID:   1,
				UserID:   1,
				Quantity: 2,
				Status:   "pending",
			},
			mockErr:        errors.New("error from service"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderService := new(mockService.MockOrderService)
			h := order.NewOrderHandler(mockOrderService)

			var body []byte
			switch v := tt.input.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
				mockOrderService.On("CreateOrder", mock.Anything).Return(tt.mockErr)
			}

			req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r := gin.Default()
			r.POST("/orders", h.CreateOrder)
			r.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestGetAllOrders(t *testing.T) {
	tests := []struct {
		name           string
		mockOrders     []*models.Order // <- Sửa ở đây
		mockErr        error
		expectedStatus int
	}{
		{
			name: "success",
			mockOrders: []*models.Order{
				{ID: 1, BookID: 1, Quantity: 2},
			},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "fail",
			mockOrders:     nil, // Tránh nil panic
			mockErr:        errors.New("db error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderService := new(mockService.MockOrderService)
			h := order.NewOrderHandler(mockOrderService)

			mockOrderService.On("GetAllOrders").Return(tt.mockOrders, tt.mockErr)

			req := httptest.NewRequest(http.MethodGet, "/orders", nil)
			w := httptest.NewRecorder()

			r := gin.Default()
			r.GET("/orders", h.GetAllOrders)
			r.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestGetByOrderID(t *testing.T) {
	tests := []struct {
		name           string
		param          string
		mockOrder      *models.Order
		mockErr        error
		expectedStatus int
	}{
		{
			name:           "valid",
			param:          "1",
			mockOrder:      &models.Order{ID: 1, Status: "processing"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid ID",
			param:          "abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "not found",
			param:          "999",
			mockOrder:      nil,
			mockErr:        errors.New("not found"),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderService := new(mockService.MockOrderService)
			h := order.NewOrderHandler(mockOrderService)

			if tt.expectedStatus != http.StatusBadRequest {
				mockOrderService.
					On("GetByOrderID", mock.Anything).
					Return(tt.mockOrder, tt.mockErr)
			}

			req := httptest.NewRequest(http.MethodGet, "/orders/"+tt.param, nil)
			w := httptest.NewRecorder()

			r := gin.Default()
			r.GET("/orders/:id", h.GetByOrderID)
			r.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
func TestDeleteByOrderID(t *testing.T) {
	tests := []struct {
		name           string
		param          string
		mockOrder      *models.Order
		mockErr        error
		expectedStatus int
	}{
		{
			name:           "Valid",
			param:          "1",
			mockOrder:      &models.Order{ID: 1},
			mockErr:        nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid param",
			param:          "abc",
			mockOrder:      nil,
			mockErr:        nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Order not found",
			param:          "999",
			mockOrder:      nil,
			mockErr:        errors.New("not found"),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "Unexpected nil order (no error)",
			param:          "2",
			mockOrder:      nil,
			mockErr:        nil,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Internal server error",
			param:          "3",
			mockOrder:      nil,
			mockErr:        errors.New("database down"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderService := new(mockService.MockOrderService)
			h := order.NewOrderHandler(mockOrderService)

			if tt.expectedStatus != http.StatusBadRequest {
				mockOrderService.
					On("DeleteByOrderID", mock.Anything).
					Return(tt.mockOrder, tt.mockErr)
			}

			req := httptest.NewRequest(http.MethodDelete, "/orders/"+tt.param, nil)
			w := httptest.NewRecorder()

			r := gin.Default()
			r.DELETE("/orders/:id", h.DeleteByOrderID)
			r.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestUpdateByOrderID(t *testing.T) {
	tests := []struct {
		name           string
		param          string
		body           interface{}
		mockReturn     *models.Order // <- sửa ở đây
		mockErr        error
		expectedStatus int
	}{
		{
			name:           "valid",
			param:          "1",
			body:           models.Order{Status: "shipped"},
			mockReturn:     &models.Order{ID: 1, Status: "shipped"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid ID",
			param:          "abc",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid JSON",
			param:          "1",
			body:           `invalid json`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "update fail",
			param:          "2",
			body:           models.Order{Status: "fail"},
			mockReturn:     nil,
			mockErr:        errors.New("fail"),
			expectedStatus: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockOrderService := new(mockService.MockOrderService)
			h := order.NewOrderHandler(mockOrderService)

			var body []byte
			switch v := tt.body.(type) {
			case string:
				body = []byte(v)
			default:
				body, _ = json.Marshal(v)
				mockOrderService.On("UpdateByOrderID", mock.Anything).Return(tt.mockReturn, tt.mockErr)
			}

			req := httptest.NewRequest(http.MethodPut, "/orders/"+tt.param, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r := gin.Default()
			r.PUT("/orders/:id", h.UpdateByOrderID)
			r.ServeHTTP(w, req)

			require.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
