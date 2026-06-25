package bootstrap

import (
	"orders-service/internal/api"
	"orders-service/internal/infra"
	"orders-service/internal/order"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InitApp acts as a manual Dependency Injection container.
// It wires up all the application dependencies (repositories, handlers)
// and returns the initialized router, keeping the main.go clean.
func InitApp(db *gorm.DB) *gin.Engine {
	// ==========================================
	// 1. Order Domain Injection
	// ==========================================
	lockService := infra.NewPGLockService(db)
	orderRepo := order.NewRepository(db)
	orderService := order.NewService(orderRepo, lockService)
	orderHandler := order.NewHandler(orderService)

	// ==========================================
	// 2. Future Domains
	// ==========================================

	// ==========================================
	// 3. Router Initialization
	// ==========================================
	return api.InitRouter(orderHandler)
}
