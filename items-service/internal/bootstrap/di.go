package bootstrap

import (
	"items-service/internal/api"
	"items-service/internal/item"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InitApp acts as a manual Dependency Injection container.
// It wires up all the application dependencies (repositories, handlers)
// and returns the initialized router, keeping the main.go clean.
func InitApp(db *gorm.DB) *gin.Engine {
	// ==========================================
	// 1. Item Domain Injection
	// ==========================================
	itemRepo := item.NewRepository(db)
	itemHandler := item.NewHandler(itemRepo)

	// ==========================================
	// 2. Future Domains (Users, Reviews, etc.)
	// ==========================================
	// userRepo := user.NewRepository(db)
	// userHandler := user.NewHandler(userRepo)

	// ==========================================
	// 3. Router Initialization
	// ==========================================
	router := api.InitRouter(itemHandler)
	return router
}
