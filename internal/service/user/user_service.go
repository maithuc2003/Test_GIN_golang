package user

import (
	"errors"

	"github.com/maithuc2003/Test_GIN_golang/internal/interfaces/repositories"
	"github.com/maithuc2003/Test_GIN_golang/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (r *UserService) GetByUsername(username string) (*models.User, error) {

	if username == "" {
		return nil, errors.New("username cannot be empty")
	}

	user, err := r.userRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func (s *UserService) LoginUser(username, password string) (*models.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}
	// hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	// fmt.Println(string(hash)) // → chuỗi mã hóa kiểu: $2a$10$...

	// fmt.Println(password)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}
