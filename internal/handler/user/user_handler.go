package user

import (
	"github.com/gin-gonic/gin"
	"github.com/maithuc2003/Test_GIN_golang/internal/interfaces/service"
	jwtutil "github.com/maithuc2003/Test_GIN_golang/pkg/jwt"
)

type UserHandler struct {
	userService service.UserServiceInterface
	JwtGenFunc  func(userID uint, username string) (string, error)
}

func NewUserHandler(userService service.UserServiceInterface) *UserHandler {
	return &UserHandler{
		userService: userService,
		JwtGenFunc: func(userID uint, username string) (string, error) {
			return jwtutil.GenerateJWT(userID, username)
		},
	}
}

func (h *UserHandler) GetByUsername(c *gin.Context) {
	username := c.Query("username")
	// username := c.Param("username")
	if username == "" {
		c.JSON(400, gin.H{"error": "Username is required"})
		return
	}
	user, err := h.userService.GetByUsername(username)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	c.JSON(200, gin.H{
		"ID":       user.ID,
		"Username": user.Username,
	})

}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	user, err := h.userService.LoginUser(req.Username, req.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := h.JwtGenFunc(user.ID, req.Username)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(200, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"token":    token,
	})
}
