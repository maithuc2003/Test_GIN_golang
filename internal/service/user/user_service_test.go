package user_test

import (
	"errors"
	"testing"

	"github.com/maithuc2003/Test_GIN_golang/internal/models"
	"github.com/maithuc2003/Test_GIN_golang/internal/service/user"
	"github.com/maithuc2003/Test_GIN_golang/internal/service/user/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestGetByUsername(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		mockReturnUser *models.User
		mockReturnErr  error
		expectedUser   *models.User
		expectedErr    string
	}{
		{
			name:         "Empty username",
			username:     "",
			expectedUser: nil,
			expectedErr:  "username cannot be empty",
		},
		{
			name:          "User not found",
			username:      "unknown_user",
			mockReturnErr: errors.New("user not found"),
			expectedUser:  nil,
			expectedErr:   "invalid username or password",
		},
		{
			name:     "Valid user",
			username: "john",
			mockReturnUser: &models.User{
				ID:       1,
				Username: "john",
				Password: "hashed_password",
			},
			mockReturnErr: nil,
			expectedUser: &models.User{
				ID:       1,
				Username: "john",
				Password: "hashed_password",
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockUserRepo)

			// Setup mock if username is not empty
			if tt.username != "" {
				mockRepo.On("GetByUsername", tt.username).
					Return(tt.mockReturnUser, tt.mockReturnErr)
			}

			service := user.NewUserService(mockRepo)
			actualUser, err := service.GetByUsername(tt.username)

			if tt.expectedErr != "" {
				assert.Nil(t, actualUser)
				assert.EqualError(t, err, tt.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, actualUser)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLoginUser(t *testing.T) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correct_password"), bcrypt.DefaultCost)

	tests := []struct {
		name           string
		username       string
		password       string
		mockUser       *models.User
		mockError      error
		expectedError  bool
		expectedResult *models.User
	}{
		{
			name:     "Valid login",
			username: "john",
			password: "correct_password",
			mockUser: &models.User{
				ID:       1,
				Username: "john",
				Password: string(hashedPassword),
			},
			mockError:      nil,
			expectedError:  false,
			expectedResult: &models.User{ID: 1, Username: "john", Password: string(hashedPassword)},
		},
		{
			name:           "User not found",
			username:       "unknown",
			password:       "any_password",
			mockUser:       nil,
			mockError:      errors.New("not found"),
			expectedError:  true,
			expectedResult: nil,
		},
		{
			name:     "Wrong password",
			username: "john",
			password: "wrong_password",
			mockUser: &models.User{
				ID:       1,
				Username: "john",
				Password: string(hashedPassword),
			},
			mockError:      nil,
			expectedError:  true,
			expectedResult: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockUserRepo)
			mockRepo.On("GetByUsername", tt.username).Return(tt.mockUser, tt.mockError)

			service := user.NewUserService(mockRepo)

			result, err := service.LoginUser(tt.username, tt.password)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.Username, result.Username)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
