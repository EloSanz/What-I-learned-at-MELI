package item

import (
	"errors"
	"net/http"

	"items-service/pkg/web"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo Repository
}

func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

// GetByID maneja la obtención de un item por ID
// GET /api/items/:id
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		web.Error(c, http.StatusBadRequest, "Missing item ID parameter")
		return
	}

	item, err := h.repo.FindByID(id)
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

type ValidateStockRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// ValidateAndReserveStock valida la existencia del item y descuenta stock
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

	item, err := h.repo.DecrementStock(id, req.Quantity)
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
