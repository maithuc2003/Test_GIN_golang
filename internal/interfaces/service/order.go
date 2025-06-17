package service

import "github.com/maithuc2003/Test_GIN_golang/internal/models"

type OrderServiceInterface interface {
	CreateOrder(order *models.Order) error
	GetAllOrders() ([]*models.Order, error)
	GetByOrderID(id int) (*models.Order, error)
	DeleteByOrderID(id int) (*models.Order, error)
	UpdateByOrderID(order *models.Order) (*models.Order, error)
}
