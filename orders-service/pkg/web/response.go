package web

import "github.com/gin-gonic/gin"

// Response representa un formato de respuesta estándar
type Response struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// JSON envía una respuesta de éxito en formato JSON
func JSON(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, Response{
		Status:  "success",
		Data:    data,
		Message: message,
	})
}

// Error envía una respuesta de error en formato JSON
func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Status:  "error",
		Message: message,
	})
}
