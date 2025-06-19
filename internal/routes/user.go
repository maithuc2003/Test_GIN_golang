package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maithuc2003/Test_GIN_golang/internal/handler/user"
	"github.com/maithuc2003/Test_GIN_golang/internal/interfaces/repositories"
	ServiceInterface "github.com/maithuc2003/Test_GIN_golang/internal/interfaces/service"
	Repo "github.com/maithuc2003/Test_GIN_golang/internal/repositories/user"
	ServiceImp "github.com/maithuc2003/Test_GIN_golang/internal/service/user"
	"gorm.io/gorm"
)

func RegisterUserRoutes(r *gin.Engine, db *gorm.DB) {
	var userRepo repositories.UserRepository = Repo.NewRepository(db)
	var userService ServiceInterface.UserServiceInterface = ServiceImp.NewUserService(userRepo)
	userHandler := user.NewUserHandler(userService)

	r.GET("/users", userHandler.GetByUsername)
	r.POST("/user/login", userHandler.LoginUser)
}
