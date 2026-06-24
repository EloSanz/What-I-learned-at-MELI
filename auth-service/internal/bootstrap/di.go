package bootstrap

import (
	"net/http"

	"auth-service/internal/auth"

	"github.com/gin-gonic/gin"
)

// InitApp acts as a manual Dependency Injection container.
// It wires up all the application dependencies and returns the initialized router.
func InitApp() *gin.Engine {
	// ==========================================
	// 1. Auth Domain Injection
	// ==========================================
	authService := auth.NewService()
	authHandler := auth.NewHandler(authService)

	// ==========================================
	// 2. Router Initialization
	// ==========================================
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP"})
	})

	apiGroup := r.Group("/api/auth")
	{
		apiGroup.POST("/login", authHandler.Login)
		apiGroup.GET("/verify", authHandler.VerifyToken)
	}

	return r
}
