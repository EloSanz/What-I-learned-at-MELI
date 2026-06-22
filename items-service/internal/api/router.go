package api

import (
	"items-service/internal/item"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitRouter configures the Gin server routes and its global middlewares
func InitRouter(itemHandler *item.Handler) *gin.Engine {
	r := gin.New()

	// Default global middlewares
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Public API group for the items domain (all under /api/items)
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
