package book

import (
	"errors"
	"fmt"

	"github.com/maithuc2003/Test_GIN_golang/internal/interfaces/repositories"
	"github.com/maithuc2003/Test_GIN_golang/internal/models"

	"gorm.io/gorm"
)

type bookRepo struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) repositories.BookRepository {
	return &bookRepo{db: db}
}

func (r *bookRepo) CreateBook(book *models.Book) error {
	return r.db.Create(book).Error
}

func (r *bookRepo) GetAllBooks() ([]models.Book, error) {
	var books []models.Book
	err := r.db.Find(&books).Error
	return books, err
}

// Lấy sách theo ID
func (r *bookRepo) GetByBookID(id int) (*models.Book, error) {
	var book models.Book
	if err := r.db.First(&book, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("book with ID %d not found", id)
		}
		return nil, err
	}
	return &book, nil
}

// Xoá sách theo ID
func (r *bookRepo) DeleteById(id int) (*models.Book, error) {
	book, err := r.GetByBookID(id)
	if err != nil {
		return nil, err
	}
	// GORM sẽ tự xử lý foreign key nếu được định nghĩa trong DB
	if err := r.db.Delete(&models.Book{}, id).Error; err != nil {
		return nil, fmt.Errorf("failed to delete book: %w", err)
	}
	return book, nil
}

// 1. Kiểm tra author tồn tại
// 2. UPDATE book

func (r *bookRepo) UpdateById(book *models.Book) (*models.Book, error) {
	// Ví dụ, dùng GORM Update:
	var count int64
	if err := r.db.Model(&models.Author{}).Where("id = ?", book.AuthorID).Count(&count).Error; err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("author not found")
	}

	result := r.db.Model(&models.Book{}).Where("id = ?", book.ID).Updates(models.Book{
		Title:     book.Title,
		AuthorID:  book.AuthorID,
		Stock:     book.Stock,
		UpdatedAt: book.UpdatedAt,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, errors.New("no book updated")
	}

	return book, nil
}
