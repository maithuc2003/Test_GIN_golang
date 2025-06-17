package repositories

import "github.com/maithuc2003/Test_GIN_golang/internal/models"

type UserRepository interface {
	GetByUsername(username string) (*models.User, error)
	LoginUser(username string, password string) (*models.User, error)
}
