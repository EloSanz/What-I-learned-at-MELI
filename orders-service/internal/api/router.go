package api

import (
	"net/http"
	"orders-service/internal/order"

	"github.com/gin-gonic/gin"
)

// InitRouter configures Gin server routes for the orders microservice
func InitRouter(orderHandler *order.Handler) *gin.Engine {
	r := gin.New()

	// Default middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Group all routes under /api/orders
	apiOrders := r.Group("/api/orders")
	{
		apiOrders.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "UP"})
		})
		apiOrders.POST("", orderHandler.CreateOrder)
		apiOrders.GET("/:id", orderHandler.GetByID)
		apiOrders.PUT("/:id/status", orderHandler.UpdateStatus)
	}

	return r
}
