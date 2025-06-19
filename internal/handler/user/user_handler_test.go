package user_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/maithuc2003/Test_GIN_golang/internal/handler/user"
	mocks "github.com/maithuc2003/Test_GIN_golang/internal/mocks/service"
	"github.com/maithuc2003/Test_GIN_golang/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetByUsernameHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		queryParam       string
		mockReturnUser   *models.User
		mockReturnErr    error
		expectedCode     int
		expectedResponse string
	}{
		{
			name:             "Missing username",
			queryParam:       "",
			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"error":"Username is required"}`,
		},
		{
			name:             "User not found",
			queryParam:       "notfound",
			mockReturnUser:   nil,
			mockReturnErr:    errors.New("user not found"),
			expectedCode:     http.StatusNotFound,
			expectedResponse: `{"error":"User not found"}`,
		},
		{
			name:       "User found successfully",
			queryParam: "john",
			mockReturnUser: &models.User{
				ID:       1,
				Username: "john",
				Password: "hashed_pw",
			},
			mockReturnErr:    nil,
			expectedCode:     http.StatusOK,
			expectedResponse: `{"ID":1,"Username":"john"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserService := new(mocks.MockUserService)
			userHandler := user.NewUserHandler(mockUserService)

			if tt.queryParam != "" {
				mockUserService.
					On("GetByUsername", tt.queryParam).
					Return(tt.mockReturnUser, tt.mockReturnErr)
			}

			router := gin.Default()
			router.GET("/users", userHandler.GetByUsername)

			req := httptest.NewRequest(http.MethodGet, "/users?username="+tt.queryParam, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedCode, resp.Code)
			assert.JSONEq(t, tt.expectedResponse, resp.Body.String())

			mockUserService.AssertExpectations(t)
		})
	}
}
func TestLoginUserHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name             string
		requestBody      map[string]string
		mockReturnUser   *models.User
		mockReturnErr    error
		expectedCode     int
		expectedResponse string
		expectedContains string // dùng khi không so sánh JSON chính xác (token thay đổi)
		mockJWTFunc      func(userID uint, username string) (string, error)
	}{
		{
			name:             "Invalid JSON",
			requestBody:      nil,
			expectedCode:     http.StatusBadRequest,
			expectedResponse: `{"error":"Invalid request body"}`,
		},
		{
			name: "Login failed - invalid credentials",
			requestBody: map[string]string{
				"username": "john",
				"password": "wrongpassword",
			},
			mockReturnUser:   nil,
			mockReturnErr:    errors.New("invalid username or password"),
			expectedCode:     http.StatusUnauthorized,
			expectedResponse: `{"error":"Invalid username or password"}`,
		},
		{
			name: "Login success - mocked JWT",
			requestBody: map[string]string{
				"username": "john",
				"password": "correctpassword",
			},
			mockReturnUser: &models.User{
				ID:       1,
				Username: "john",
			},
			mockReturnErr:    nil,
			expectedCode:     http.StatusOK,
			expectedResponse: `{"id":1,"username":"john","token":"mocked.token.jwt"}`,
			mockJWTFunc: func(userID uint, username string) (string, error) {
				return "mocked.token.jwt", nil
			},
		},
		{
			name: "Login failed - JWT generation error",
			requestBody: map[string]string{
				"username": "john",
				"password": "correctpassword",
			},
			mockReturnUser: &models.User{
				ID:       1,
				Username: "john",
			},
			mockReturnErr:    nil,
			expectedCode:     http.StatusInternalServerError,
			expectedResponse: `{"error":"Failed to generate token"}`,
			mockJWTFunc: func(userID uint, username string) (string, error) {
				return "", errors.New("token generation error")
			},
		},
		{
			name: "Login success - real JWT",
			requestBody: map[string]string{
				"username": "realuser",
				"password": "realpass",
			},
			mockReturnUser: &models.User{
				ID:       2,
				Username: "realuser",
			},
			mockReturnErr:    nil,
			expectedCode:     http.StatusOK,
			expectedContains: `"id":2,"username":"realuser","token":"`, // chỉ cần chứa token là ok
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserService := new(mocks.MockUserService)
			userHandler := user.NewUserHandler(mockUserService)

			// Nếu có mock JWT thì gán lại
			if tt.mockJWTFunc != nil {
				userHandler.JwtGenFunc = func(userID uint, username string) (string, error) {
					return tt.mockJWTFunc(userID, username)
				}
			}

			// Gán mock LoginUser nếu requestBody hợp lệ
			if tt.requestBody != nil {
				mockUserService.
					On("LoginUser", tt.requestBody["username"], tt.requestBody["password"]).
					Return(tt.mockReturnUser, tt.mockReturnErr)
			}

			router := gin.Default()
			router.POST("/login", userHandler.LoginUser)

			var reqBodyBytes []byte
			if tt.requestBody != nil {
				reqBodyBytes, _ = json.Marshal(tt.requestBody)
			} else {
				reqBodyBytes = []byte("invalid json") // simulate invalid json
			}

			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(reqBodyBytes))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedCode, resp.Code)

			if tt.expectedResponse != "" {
				assert.JSONEq(t, tt.expectedResponse, resp.Body.String())
			}

			if tt.name == "Login success - real JWT" {
				var respBody struct {
					ID       int    `json:"id"`
					Username string `json:"username"`
					Token    string `json:"token"`
				}
				err := json.Unmarshal(resp.Body.Bytes(), &respBody)
				assert.NoError(t, err)

				assert.Equal(t, 2, respBody.ID)
				assert.Equal(t, "realuser", respBody.Username)
				assert.NotEmpty(t, respBody.Token)
			} else if tt.expectedResponse != "" {
				assert.JSONEq(t, tt.expectedResponse, resp.Body.String())
			}

			mockUserService.AssertExpectations(t)
		})
	}
}
