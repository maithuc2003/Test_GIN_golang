package service

import "github.com/maithuc2003/Test_GIN_golang/internal/models"


type AuthorServiceInterface interface {
	CreateAuthor(author *models.Author) error
	GetAllAuthors() ([]*models.Author, error)
	GetByAuthorID(id int) (*models.Author, error)
	DeleteById(id int) (*models.Author, error)
	UpdateById(author *models.Author) (*models.Author, error)
}
