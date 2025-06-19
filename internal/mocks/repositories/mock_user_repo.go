package mocks

import (
	"github.com/maithuc2003/Test_GIN_golang/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) GetByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	user, _ := args.Get(0).(*models.User)
	return user, args.Error(1)
}

func (m *MockUserRepo) LoginUser(username string, password string) (*models.User, error) {
	args := m.Called(username, password)
	user, ok := args.Get(0).(*models.User)
	if !ok {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}
