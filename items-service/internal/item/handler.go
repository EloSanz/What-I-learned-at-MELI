package item

import (
	"errors"
	"net/http"

	"items-service/pkg/web"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetByID handles retrieving an item by ID
// GET /api/items/:id
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		web.Error(c, http.StatusBadRequest, "Missing item ID parameter")
		return
	}

	item, err := h.service.GetByID(id)
	if err != nil {
		if errors.Is(err, ErrItemNotFound) {
			web.Error(c, http.StatusNotFound, "Item not found")
			return
		}
		web.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.JSON(c, http.StatusOK, item, "Item retrieved successfully")
}

// ValidateAndReserveStock validates item existence and reserves stock
// POST /api/items/:id/validate
func (h *Handler) ValidateAndReserveStock(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		web.Error(c, http.StatusBadRequest, "Missing item ID parameter")
		return
	}

	var req ValidateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		web.Error(c, http.StatusBadRequest, "Invalid request payload. Quantity must be 1 or greater.")
		return
	}

	item, err := h.service.ValidateAndReserveStock(id, req.Quantity)
	if err != nil {
		if errors.Is(err, ErrItemNotFound) {
			web.Error(c, http.StatusNotFound, "Item not found")
			return
		}
		if errors.Is(err, ErrOutOfStock) {
			web.Error(c, http.StatusUnprocessableEntity, "Insufficient stock available for this purchase")
			return
		}
		web.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	web.JSON(c, http.StatusOK, item, "Stock validated and decremented successfully")
}
