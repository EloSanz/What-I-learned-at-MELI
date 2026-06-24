package auth

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests for the authentication domain
type Handler struct {
	service Service
}

// NewHandler creates a new auth handler
func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// LoginRequest defines the expected JSON payload for login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login authenticates a user and returns a JWT
// POST /api/auth/login
func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid request payload"})
		return
	}

	// Hardcoded credentials for demonstration purposes.
	// In a real app, you would query a Users database here.
	if req.Username != "admin" || req.Password != "admin" {
		slog.Warn("Failed login attempt", "username", req.Username)
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid credentials"})
		return
	}

	// Generate JWT for the user
	token, err := h.service.GenerateToken(req.Username)
	if err != nil {
		slog.Error("Failed to generate token", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Could not generate token"})
		return
	}

	slog.Info("User logged in successfully", "username", req.Username)
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"token": token,
		},
		"message": "Login successful",
	})
}

// VerifyToken checks the Authorization header and validates the JWT
// GET /api/auth/verify
func (h *Handler) VerifyToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		slog.Warn("Authentication failed: Missing Authorization header")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
		return
	}

	// Expected format: "Bearer <token>"
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		slog.Warn("Authentication failed: Invalid token format")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format, expected Bearer <token>"})
		return
	}

	tokenString := parts[1]

	// Validate the JWT cryptographically
	_, err := h.service.ValidateToken(tokenString)
	if err != nil {
		slog.Warn("Authentication failed: Invalid JWT", "error", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	// Token is valid. Nginx only cares about the 200 OK status.
	// In a more advanced setup, we could pass the decoded user ID back in headers.
	c.JSON(http.StatusOK, gin.H{"status": "authenticated"})
}
