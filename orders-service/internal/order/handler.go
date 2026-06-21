package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"orders-service/pkg/web"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	repo            Repository
	itemsServiceURL string
	httpClient      *http.Client
}

func NewHandler(repo Repository) *Handler {
	itemsURL := os.Getenv("ITEMS_SERVICE_URL")
	if itemsURL == "" {
		itemsURL = "http://localhost:8081"
	}

	return &Handler{
		repo:            repo,
		itemsServiceURL: itemsURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// ItemResponseData mapea la sección de datos recibida desde items-service
type ItemResponseData struct {
	ID    string  `json:"id"`
	Title string  `json:"title"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

// ItemsServiceResponse mapea la respuesta completa de items-service
type ItemsServiceResponse struct {
	Status  string           `json:"status"`
	Data    ItemResponseData `json:"data"`
	Message string           `json:"message"`
}

type CreateOrderRequest struct {
	UserID   string  `json:"user_id" binding:"required"`
	ItemID   string  `json:"item_id" binding:"required"`
	Quantity int     `json:"quantity" binding:"required,min=1"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
	Address  string  `json:"address" binding:"required"`
}

// CreateOrder maneja la creación de una orden y valida el item/stock con items-service
// POST /api/orders
func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		web.Error(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	// 1. Validar existencia y stock llamando al microservicio de items
	slog.Info("Validating Item in items-service", "item_id", req.ItemID, "quantity", req.Quantity)
	itemDetail, err := h.checkItemAvailability(req.ItemID)
	if err != nil {
		slog.Error("Item validation failed", "error", err)
		web.Error(c, http.StatusBadRequest, "Item validation failed: "+err.Error())
		return
	}

	// 2. Validar stock localmente antes de crear la orden
	if itemDetail.Stock < req.Quantity {
		web.Error(c, http.StatusUnprocessableEntity, fmt.Sprintf("Insufficient stock. Available: %d, Requested: %d", itemDetail.Stock, req.Quantity))
		return
	}

	// 3. Crear el modelo de orden
	ord := Order{
		ID:       GenerateUUID(),
		UserID:   req.UserID,
		ItemID:   req.ItemID,
		Quantity: req.Quantity,
		Amount:   req.Amount,
		Address:  req.Address,
		Status:   StatusReadyToProcess, // Seteamos el estado para el orquestador
	}

	// 4. Persistir en la base de datos
	if err := h.repo.Create(&ord); err != nil {
		web.Error(c, http.StatusInternalServerError, "Failed to persist order: "+err.Error())
		return
	}

	slog.Info("Order successfully created", "order_id", ord.ID, "status", ord.Status)
	web.JSON(c, http.StatusCreated, ord, "Order generated successfully")
}

// GetByID busca una orden por ID
// GET /api/orders/:id
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		web.Error(c, http.StatusBadRequest, "Missing order ID parameter")
		return
	}

	ord, err := h.repo.FindByID(id)
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

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// UpdateStatus permite cambiar el estado de la orden (usado por el orquestador)
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

	// Validar que sea un estado admitido
	if req.Status != StatusPending && req.Status != StatusReadyToProcess && req.Status != StatusCompleted && req.Status != StatusFailed {
		web.Error(c, http.StatusBadRequest, "Invalid order status value")
		return
	}

	if err := h.repo.UpdateStatus(id, req.Status); err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			web.Error(c, http.StatusNotFound, "Order not found")
			return
		}
		web.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	slog.Info("Order status updated", "order_id", id, "new_status", req.Status)
	web.JSON(c, http.StatusOK, gin.H{"id": id, "status": req.Status}, "Order status updated successfully")
}

// Helper para llamar por HTTP al microservicio de items
func (h *Handler) checkItemAvailability(itemID string) (*ItemResponseData, error) {
	url := fmt.Sprintf("%s/api/items/%s", h.itemsServiceURL, itemID)
	resp, err := h.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("items-service is unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("the requested item does not exist")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("items-service returned status code %d", resp.StatusCode)
	}

	var response ItemsServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response payload: %w", err)
	}

	if response.Status != "success" {
		return nil, errors.New("failed response status from items-service")
	}

	return &response.Data, nil
}
