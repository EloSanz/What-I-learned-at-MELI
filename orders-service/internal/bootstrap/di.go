package bootstrap

import (
	"orders-service/internal/api"
	"orders-service/internal/broker"
	"orders-service/internal/order"
	"os"

	sdkbroker "github.com/user/meli-sdk/broker"
	"github.com/user/meli-sdk/lock"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// InitApp acts as a manual Dependency Injection container.
// It wires up all the application dependencies (repositories, handlers)
// and returns the initialized router, keeping the main.go clean.
func InitApp(db *gorm.DB) *gin.Engine {
	// ==========================================
	// 0. RabbitMQ Connection
	// ==========================================
	rabbitURL := os.Getenv("RABBITMQ_URL")
	if rabbitURL == "" {
		rabbitURL = "amqp://user:password@localhost:5672/"
	}
	rabbitMQ, err := sdkbroker.NewRabbitMQ(rabbitURL)
	if err != nil {
		panic("Failed to connect to RabbitMQ: " + err.Error())
	}

	// ==========================================
	// 1. Order Domain Injection
	// ==========================================
	lockService := lock.NewPGLockService(db)
	orderRepo := order.NewRepository(db)
	orderService := order.NewService(orderRepo, lockService, rabbitMQ)
	orderHandler := order.NewHandler(orderService)

	// ==========================================
	// 2. Start Async Consumers
	// ==========================================
	orderConsumer := broker.NewConsumer(rabbitMQ, orderService)
	orderConsumer.Start()

	// ==========================================
	// 2. Future Domains
	// ==========================================

	// ==========================================
	// 3. Router Initialization
	// ==========================================
	return api.InitRouter(orderHandler)
}
