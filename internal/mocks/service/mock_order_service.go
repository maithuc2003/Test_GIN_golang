package mocks

import (
	"github.com/maithuc2003/Test_GIN_golang/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrder(order *models.Order) error {
	args := m.Called(order)
	return args.Error(0)
}

func (m *MockOrderService) GetAllOrders() ([]*models.Order, error) {
	args := m.Called()
	return args.Get(0).([]*models.Order), args.Error(1)
}


func (m *MockOrderService) GetByOrderID(id int) (*models.Order, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}

func (m *MockOrderService) DeleteByOrderID(id int) (*models.Order, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Order), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockOrderService) UpdateByOrderID(order *models.Order) (*models.Order, error) {
	args := m.Called(order)
	if args.Get(0) != nil {
		return args.Get(0).(*models.Order), args.Error(1)
	}
	return nil, args.Error(1)
}
