package order_test

import (
	"errors"
	"testing"
	"time"

	mocks "github.com/maithuc2003/Test_GIN_golang/internal/mocks/repositories"
	"github.com/maithuc2003/Test_GIN_golang/internal/models"
	"github.com/maithuc2003/Test_GIN_golang/internal/service/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateOrder(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	service := order.NewOrderService(mockRepo)

	tests := []struct {
		name        string
		input       *models.Order
		mockError   error
		expectError bool
	}{
		{
			name:        "nil order",
			input:       nil,
			expectError: true,
		},
		{
			name:        "invalid book ID",
			input:       &models.Order{BookID: 0, UserID: 1, Quantity: 1, Status: "pending"},
			expectError: true,
		},
		{
			name:        "invalid user ID",
			input:       &models.Order{BookID: 1, UserID: 0, Quantity: 1, Status: "pending"},
			expectError: true,
		},
		{
			name:        "quantity <= 0",
			input:       &models.Order{BookID: 1, UserID: 1, Quantity: 0, Status: "pending"},
			expectError: true,
		},
		{
			name:        "status empty",
			input:       &models.Order{BookID: 1, UserID: 1, Quantity: 1, Status: ""},
			expectError: true,
		},
		{
			name:        "valid order",
			input:       &models.Order{BookID: 1, UserID: 1, Quantity: 1, Status: "confirmed"},
			mockError:   nil,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.input != nil && !tt.expectError {
				mockRepo.On("Create", mock.AnythingOfType("*models.Order")).Return(tt.mockError).Once()
			}
			err := service.CreateOrder(tt.input)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
func TestOrderService_GetAllOrders(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	s := order.NewOrderService(mockRepo)

	tests := []struct {
		name        string
		mockOrders  []*models.Order
		mockError   error
		expectedErr string
		expectedLen int
	}{
		{
			name:        "success",
			mockOrders:  []*models.Order{{ID: 1}, {ID: 2}},
			mockError:   nil,
			expectedErr: "",
			expectedLen: 2,
		},
		{
			name:        "return error when no orders found",
			mockOrders:  []*models.Order{}, // empty slice
			mockError:   nil,
			expectedErr: "no orders found",
			expectedLen: 0,
		},
		{
			name:        "repository error",
			mockOrders:  nil,
			mockError:   errors.New("db error"),
			expectedErr: "db error",
			expectedLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil // reset mock
			mockRepo.On("GetAllOrders").Return(tt.mockOrders, tt.mockError).Once()

			result, err := s.GetAllOrders()

			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr)
				assert.Nil(t, result)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.expectedLen)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestOrderService_GetByOrderID(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	s := order.NewOrderService(mockRepo)

	tests := []struct {
		name        string
		id          int
		mockOrder   *models.Order
		mockError   error
		expectedErr string
	}{
		{
			name:        "invalid ID",
			id:          -1,
			expectedErr: "invalid order ID",
		},
		{
			name:        "success",
			id:          1,
			mockOrder:   &models.Order{ID: 1},
			mockError:   nil,
			expectedErr: "",
		},
		{
			name:        "repo error",
			id:          2,
			mockError:   errors.New("not found"),
			expectedErr: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.id > 0 {
				mockRepo.On("GetByOrderID", uint(tt.id)).Return(tt.mockOrder, tt.mockError).Once()
			}

			result, err := s.GetByOrderID(tt.id)

			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.mockOrder, result)
			}
		})
	}
}

func TestOrderService_DeleteByOrderID(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	s := order.NewOrderService(mockRepo)

	tests := []struct {
		name        string
		id          int
		mockOrder   *models.Order
		mockError   error
		expectedErr string
	}{
		{
			name:        "invalid ID",
			id:          0,
			expectedErr: "invalid order ID",
		},
		{
			name:        "success",
			id:          1,
			mockOrder:   &models.Order{ID: 1},
			mockError:   nil,
			expectedErr: "",
		},
		{
			name:        "repo error",
			id:          2,
			mockError:   errors.New("delete error"),
			expectedErr: "delete error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.id > 0 {
				mockRepo.On("DeleteByOrderID", uint(tt.id)).Return(tt.mockOrder, tt.mockError).Once()
			}

			result, err := s.DeleteByOrderID(tt.id)

			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.mockOrder, result)
			}
		})
	}
}

func TestOrderService_UpdateByOrderID(t *testing.T) {
	mockRepo := new(mocks.MockOrderRepository)
	s := order.NewOrderService(mockRepo)

	// now := time.Now()

	tests := []struct {
		name        string
		order       *models.Order
		mockReturn  *models.Order
		mockError   error
		expectedErr string
	}{
		{
			name:        "nil order",
			order:       nil,
			expectedErr: "order is nil",
		},
		{
			name: "invalid order ID",
			order: &models.Order{
				ID: 0, BookID: 1, UserID: 1, Quantity: 1, Status: "ok",
			},
			expectedErr: "invalid order ID",
		},
		{
			name: "invalid book ID",
			order: &models.Order{
				ID: 1, BookID: 0, UserID: 1, Quantity: 1, Status: "ok",
			},
			expectedErr: "invalid book ID",
		},
		{
			name: "invalid user ID",
			order: &models.Order{
				ID: 1, BookID: 1, UserID: 0, Quantity: 1, Status: "ok",
			},
			expectedErr: "invalid user ID",
		},
		{
			name: "invalid quantity",
			order: &models.Order{
				ID: 1, BookID: 1, UserID: 1, Quantity: 0, Status: "ok",
			},
			expectedErr: "quantity must be greater than zero",
		},
		{
			name: "empty status",
			order: &models.Order{
				ID: 1, BookID: 1, UserID: 1, Quantity: 1, Status: " ",
			},
			expectedErr: "status is required",
		},
		{
			name: "success",
			order: &models.Order{
				ID: 1, BookID: 1, UserID: 1, Quantity: 1, Status: "pending",
			},
			mockReturn:  &models.Order{ID: 1},
			mockError:   nil,
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedErr == "" {
				mockRepo.On("UpdateByOrderID", tt.order).Return(tt.mockReturn, tt.mockError).Once()
			}

			result, err := s.UpdateByOrderID(tt.order)

			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.mockReturn, result)
				assert.WithinDuration(t, time.Now(), tt.order.UpdatedAt, time.Second)
			}
		})
	}
}
