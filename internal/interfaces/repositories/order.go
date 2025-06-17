package repositories

import "github.com/maithuc2003/Test_GIN_golang/internal/models"

type OrderRepositoryInterface interface {
	GetByOrderID(id uint) (*models.Order, error)
	GetAllOrders() ([]*models.Order, error)
	UpdateByOrderID(order *models.Order) (*models.Order, error)
	DeleteByOrderID(id uint) (*models.Order, error)
	Create(order *models.Order) error
}
