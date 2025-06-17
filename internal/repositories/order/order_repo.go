package repositories

import (
	"errors"
	"fmt"

	"github.com/maithuc2003/Test_GIN_golang/internal/interfaces/repositories"
	"github.com/maithuc2003/Test_GIN_golang/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type orderRepo struct {
	db *gorm.DB
}

// Đảm bảo đúng interface name, tên đúng OrderRepositoryInterface
func NewOrderRepo(db *gorm.DB) repositories.OrderRepositoryInterface {
	return &orderRepo{db: db}
}

// Tạo đơn hàng, giảm stock sách trong transaction có khóa dòng (pessimistic lock)
func (r *orderRepo) Create(order *models.Order) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var book models.Book
		// Khóa bản ghi sách đang xử lý để tránh race condition
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&book, order.BookID).Error; err != nil {
			return fmt.Errorf("book not found: %w", err)
		}

		if book.Stock < order.Quantity {
			return fmt.Errorf("not enough stock available")
		}

		if err := tx.Create(order).Error; err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}

		if err := tx.Model(&book).Update("stock", book.Stock-order.Quantity).Error; err != nil {
			return fmt.Errorf("failed to update book stock: %w", err)
		}

		return nil
	})
}

// Lấy tất cả đơn hàng
func (r *orderRepo) GetAllOrders() ([]*models.Order, error) {
	var orders []*models.Order
	if err := r.db.Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

// Lấy đơn hàng theo ID
func (r *orderRepo) GetByOrderID(id uint) (*models.Order, error) {
	var order models.Order
	if err := r.db.First(&order, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order with ID %d not found", id)
		}
		return nil, err
	}
	return &order, nil
}

// Xóa đơn hàng theo ID và trả về đơn hàng đã xóa
func (r *orderRepo) DeleteByOrderID(id uint) (*models.Order, error) {
	order, err := r.GetByOrderID(id)
	if err != nil {
		return nil, err
	}

	if err := r.db.Delete(&models.Order{}, id).Error; err != nil {
		return nil, fmt.Errorf("failed to delete order: %w", err)
	}

	return order, nil
}

// Cập nhật đơn hàng theo ID, trả về đơn hàng đã cập nhật
func (r *orderRepo) UpdateByOrderID(order *models.Order) (*models.Order, error) {
	result := r.db.Model(&models.Order{}).
		Where("id = ?", order.ID).
		Updates(map[string]interface{}{
			"book_id":  order.BookID,
			"user_id":  order.UserID,
			"quantity": order.Quantity,
			"status":   order.Status,
		})

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("no order updated with id %d", order.ID)
	}

	updatedOrder, err := r.GetByOrderID(order.ID)
	if err != nil {
		return nil, err
	}

	return updatedOrder, nil
}
