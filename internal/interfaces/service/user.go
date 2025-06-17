package service

import "github.com/maithuc2003/Test_GIN_golang/internal/models"

type UserServiceInterface interface {
	GetByUsername(username string) (*models.User, error)
	LoginUser(username string, password string) (*models.User, error)
}
