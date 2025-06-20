package repositories

import "github.com/maithuc2003/Test_GIN_golang/internal/models"

type AuthorRepositoriesInterface interface {
	GetByAuthorID(id int) (*models.Author, error)
	GetAllAuthors() ([]*models.Author, error)
	CreateAuthor(author *models.Author) error
	UpdateById(author *models.Author) (*models.Author, error)
	DeleteById(id int) (*models.Author, error)
}
