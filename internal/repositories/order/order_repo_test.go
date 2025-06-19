package order_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/maithuc2003/Test_GIN_golang/internal/models"
	"github.com/maithuc2003/Test_GIN_golang/internal/repositories/order"
	sqlitedriver "gorm.io/driver/sqlite"

	_ "modernc.org/sqlite"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := fmt.Sprintf("file:testdb_%d?mode=memory&cache=shared", time.Now().UnixNano())

	db, err := gorm.Open(sqlitedriver.New(sqlitedriver.Config{
		DSN:        dsn,
		DriverName: "sqlite",
	}), &gorm.Config{})

	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.Book{}, &models.Order{}))

	return db
}

func seedBook(t *testing.T, db *gorm.DB, stock int) models.Book {
	book := models.Book{
		Title:    "Test Book",
		Stock:    stock,
		AuthorID: 1,
	}
	require.NoError(t, db.Create(&book).Error)
	return book
}

func TestOrderRepo_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := order.NewOrderRepo(db)
	book := seedBook(t, db, 10)

	tests := []struct {
		name        string
		order       models.Order
		expectedErr string
	}{
		{
			name: "success",
			order: models.Order{
				BookID:   uint(book.ID),
				UserID:   1,
				Quantity: 2,
			},
		},
		{
			name: "not enough stock",
			order: models.Order{
				BookID:   book.ID,
				UserID:   2,
				Quantity: 999,
			},
			expectedErr: "not enough stock",
		},
		{
			name: "book not found",
			order: models.Order{
				BookID:   9999,
				UserID:   3,
				Quantity: 1,
			},
			expectedErr: "book not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Create(&tt.order)
			if tt.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestOrderRepo_GetByOrderID(t *testing.T) {
	db := setupTestDB(t)
	repo := order.NewOrderRepo(db)
	book := seedBook(t, db, 5)

	order := models.Order{BookID: book.ID, UserID: 1, Quantity: 1}
	require.NoError(t, repo.Create(&order))

	tests := []struct {
		name        string
		id          uint
		expectFound bool
	}{
		{
			name:        "found",
			id:          order.ID,
			expectFound: true,
		},
		{
			name:        "not found",
			id:          9999,
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.GetByOrderID(tt.id)
			if tt.expectFound {
				require.NoError(t, err)
				require.Equal(t, order.BookID, got.BookID)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestOrderRepo_GetAllOrders(t *testing.T) {
	db := setupTestDB(t)
	repo := order.NewOrderRepo(db)
	book := seedBook(t, db, 5)

	orders := []models.Order{
		{BookID: book.ID, UserID: 1, Quantity: 1},
		{BookID: book.ID, UserID: 2, Quantity: 2},
	}

	for _, o := range orders {
		require.NoError(t, repo.Create(&o))
	}

	t.Run("get all orders", func(t *testing.T) {
		results, err := repo.GetAllOrders()
		require.NoError(t, err)
		require.Len(t, results, len(orders))
	})
}

func TestOrderRepo_DeleteByOrderID(t *testing.T) {
	db := setupTestDB(t)
	repo := order.NewOrderRepo(db)
	book := seedBook(t, db, 5)

	order := models.Order{BookID: book.ID, UserID: 1, Quantity: 1}
	require.NoError(t, repo.Create(&order))

	tests := []struct {
		name        string
		id          uint
		expectErr   bool
		expectFound bool
	}{
		{
			name:        "delete success",
			id:          order.ID,
			expectErr:   false,
			expectFound: false,
		},
		{
			name:        "delete non-existent",
			id:          9999,
			expectErr:   true,
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			deleted, err := repo.DeleteByOrderID(tt.id)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.id, deleted.ID)

				_, err := repo.GetByOrderID(tt.id)
				require.Error(t, err)
			}
		})
	}
}

func TestOrderRepo_UpdateByOrderID(t *testing.T) {
	db := setupTestDB(t)
	repo := order.NewOrderRepo(db)
	book := seedBook(t, db, 5)

	order := models.Order{BookID: book.ID, UserID: 1, Quantity: 1, Status: "Pending"}
	require.NoError(t, repo.Create(&order))

	tests := []struct {
		name        string
		update      models.Order
		expectedErr string
	}{
		{
			name: "update success",
			update: models.Order{
				ID:       order.ID,
				BookID:   book.ID,
				UserID:   1,
				Quantity: 1,
				Status:   "Completed",
			},
		},
		{
			name: "update non-existent",
			update: models.Order{
				ID:       9999,
				BookID:   book.ID,
				UserID:   1,
				Quantity: 1,
				Status:   "Shipped",
			},
			expectedErr: "no order updated",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := repo.UpdateByOrderID(&tt.update)
			if tt.expectedErr != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.expectedErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.update.Status, updated.Status)
			}
		})
	}
}
