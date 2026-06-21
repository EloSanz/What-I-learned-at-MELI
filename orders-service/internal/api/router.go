package api

import (
	"net/http"
	"orders-service/internal/order"

	"github.com/gin-gonic/gin"
)

// InitRouter configura las rutas del servidor Gin para el microservicio de órdenes
func InitRouter(orderHandler *order.Handler) *gin.Engine {
	r := gin.New()

	// Middlewares por defecto
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Agrupamos todas las rutas bajo /api/orders
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
