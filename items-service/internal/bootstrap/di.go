package bootstrap

import (
	"items-service/internal/api"
	"items-service/internal/broker"
	"items-service/internal/item"
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
	// 1. Item Domain Injection
	// ==========================================
	lockService := lock.NewPGLockService(db)
	itemRepo := item.NewRepository(db)
	itemService := item.NewService(itemRepo, lockService)
	itemHandler := item.NewHandler(itemService)

	// ==========================================
	// 2. Start Async Consumers
	// ==========================================
	itemConsumer := broker.NewConsumer(rabbitMQ, itemService)
	itemConsumer.Start()

	// ==========================================
	// 3. Future Domains (Users, Reviews, etc.)
	// ==========================================
	// userRepo := user.NewRepository(db)
	// userHandler := user.NewHandler(userRepo)

	// ==========================================
	// 4. Router Initialization
	// ==========================================
	router := api.InitRouter(itemHandler)
	return router
}
