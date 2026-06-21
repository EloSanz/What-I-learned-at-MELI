package api

import (
	"items-service/internal/item"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitRouter configura las rutas del servidor Gin y sus middlewares globales
func InitRouter(itemHandler *item.Handler) *gin.Engine {
	r := gin.New()

	// Middlewares globales por defecto
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Grupo de APIs públicas para el dominio de items (todas bajo /api/items)
	apiItems := r.Group("/api/items")
	{
		apiItems.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "UP"})
		})
		apiItems.GET("/:id", itemHandler.GetByID)
		apiItems.POST("/:id/validate", itemHandler.ValidateAndReserveStock)
	}

	return r
}
