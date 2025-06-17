package author

import (
	"fmt"

	"github.com/maithuc2003/Test_GIN_golang/internal/interfaces/repositories"
	"github.com/maithuc2003/Test_GIN_golang/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type authorRepo struct {
	db *gorm.DB
}

func NewAuthorRepo(db *gorm.DB) repositories.AuthorRepositoriesInterface {
	return &authorRepo{db: db}
}

func (r *authorRepo) GetAllAuthors() ([]*models.Author, error) {
	var authors []*models.Author
	if err := r.db.Find(&authors).Error; err != nil {
		return nil, fmt.Errorf("failed to query authors: %w", err)
	}
	return authors, nil
}

func (r *authorRepo) GetByAuthorID(id int) (*models.Author, error) {
	var author models.Author
	if err := r.db.First(&author, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("author with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch author: %w", err)
	}
	return &author, nil
}

func (r *authorRepo) CreateAuthor(author *models.Author) error {
	if err := r.db.Create(author).Error; err != nil {
		return fmt.Errorf("failed to create author: %w", err)
	}
	return nil
}

func (r *authorRepo) DeleteById(id int) (*models.Author, error) {
	author, err := r.GetByAuthorID(id)
	if err != nil {
		return nil, err
	}
	if err := r.db.Delete(&models.Author{}, id).Error; err != nil {
		// Check if it's a foreign key constraint
		return nil, fmt.Errorf("failed to delete author: %w", err)
	}
	return author, nil
}

func (r *authorRepo) UpdateById(author *models.Author) (*models.Author, error) {
	var existing models.Author
	if err := r.db.First(&existing, author.ID).Error; err != nil {
		return nil, fmt.Errorf("author_id %d does not exist", author.ID)
	}

	// Cập nhật thông tin
	if err := r.db.Model(&existing).
		Clauses(clause.Returning{}).
		Updates(map[string]interface{}{
			"name":        author.Name,
			"nationality": author.Nationality,
			"updated_at":  author.UpdatedAt,
		}).Error; err != nil {
		return nil, fmt.Errorf("failed to update author: %w", err)
	}
	return &existing, nil
}
