package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maithuc2003/Test_GIN_golang/internal/handler/order"
	RepInterface "github.com/maithuc2003/Test_GIN_golang/internal/interfaces/repositories"
	ServiceInterface "github.com/maithuc2003/Test_GIN_golang/internal/interfaces/service"
	"github.com/maithuc2003/Test_GIN_golang/internal/middleware"
	Repo "github.com/maithuc2003/Test_GIN_golang/internal/repositories/order"
	ServiceImp "github.com/maithuc2003/Test_GIN_golang/internal/service/order"
	"gorm.io/gorm"
)

func RegisterOrderRoutes(r *gin.Engine, db *gorm.DB) {
	var orderRepo RepInterface.OrderRepositoryInterface = Repo.NewOrderRepo(db)
	var orderService ServiceInterface.OrderServiceInterface = ServiceImp.NewOrderService(orderRepo)
	orderHandler := order.NewOrderHandler(orderService)

	// Public routes
	r.GET("/orders", orderHandler.GetAllOrders)
	r.GET("/orders/:id", orderHandler.GetByOrderID)

	// Authenticated (admin) routes
	auth := r.Group("/")
	auth.Use(middleware.AuthMiddleware())
	auth.POST("/orders", orderHandler.CreateOrder)
	auth.PUT("/orders/:id", orderHandler.UpdateByOrderID)
	auth.DELETE("/orders/:id", orderHandler.DeleteByOrderID)
}
