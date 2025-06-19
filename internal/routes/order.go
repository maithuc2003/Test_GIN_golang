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

	authorRoutes := r.Group("/orders")
	{
		authorRoutes.GET("", orderHandler.GetAllOrders)
		authorRoutes.GET("/:id", orderHandler.GetByOrderID)
	}

	// Protected author routes
	auth := r.Group("/orders", middleware.AuthMiddleware())
	{
		auth.POST("/add", middleware.RBACMiddleware("order/create"), orderHandler.CreateOrder)
		auth.PUT("/:id", middleware.RBACMiddleware("order/update"), orderHandler.UpdateByOrderID)
		auth.DELETE("/:id", middleware.RBACMiddleware("order/delete"), orderHandler.DeleteByOrderID)
	}

}
