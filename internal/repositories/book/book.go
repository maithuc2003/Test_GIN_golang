package repositories

import (
	"github.com/maithuc2003/GIN_golang_framework/internal/interfaces/repositories"
	"github.com/maithuc2003/GIN_golang_framework/internal/models"

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
