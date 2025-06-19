package mocks

import (
	"github.com/maithuc2003/Test_GIN_golang/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockOrderRepository mocks the OrderRepositoryInterface
type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderRepository) GetAllOrders() ([]*models.Order, error) {
	args := m.Called()
	return args.Get(0).([]*models.Order), args.Error(1)
}

func (m *MockOrderRepository) GetByOrderID(id uint) (*models.Order, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) UpdateByOrderID(order *models.Order) (*models.Order, error) {
	args := m.Called(order)
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderRepository) DeleteByOrderID(id uint) (*models.Order, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Order), args.Error(1)
}
