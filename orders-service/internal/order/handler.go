package order

import (
	"errors"
	"log/slog"
	"net/http"

	"orders-service/pkg/web"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

// CreateOrder handles order creation and validates item/stock with items-service
// POST /api/orders
func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		web.Error(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	ord, err := h.service.CreateOrder(req)
	if err != nil {
		if errors.Is(err, ErrInsufficientStock) || errors.Is(err, ErrItemValidationFailed) {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		web.Error(c, http.StatusInternalServerError, "Failed to create order: "+err.Error())
		return
	}

	web.JSON(c, http.StatusCreated, ord, "Order generated successfully")
}

// GetByID retrieves an order by ID
// GET /api/orders/:id
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		web.Error(c, http.StatusBadRequest, "Missing order ID parameter")
		return
	}

	ord, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			web.Error(c, http.StatusNotFound, "Order not found")
			return
		}
		web.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.JSON(c, http.StatusOK, ord, "Order retrieved successfully")
}

// UpdateStatus allows changing the order status (used by the orchestrator)
// PUT /api/orders/:id/status
func (h *Handler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		web.Error(c, http.StatusBadRequest, "Missing order ID parameter")
		return
	}

	var req UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		web.Error(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	if err := h.service.UpdateStatus(id, req.Status); err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			web.Error(c, http.StatusNotFound, "Order not found")
			return
		}
		if errors.Is(err, ErrInvalidOrderStatus) {
			web.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		web.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	slog.Info("Order status updated", "order_id", id, "new_status", req.Status)
	web.JSON(c, http.StatusOK, gin.H{"id": id, "status": req.Status}, "Order status updated successfully")
}
