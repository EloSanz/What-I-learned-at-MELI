package order

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

var (
	ErrItemValidationFailed = errors.New("item validation failed")
	ErrInsufficientStock    = errors.New("insufficient stock")
	ErrInvalidOrderStatus   = errors.New("invalid order status value")
)

type Service interface {
	CreateOrder(req CreateOrderRequest) (*Order, error)
	GetByID(id string) (*Order, error)
	UpdateStatus(id, status string) error
}

type service struct {
	repo            Repository
	itemsServiceURL string
	httpClient      *http.Client
}

func NewService(repo Repository) Service {
	itemsURL := os.Getenv("ITEMS_SERVICE_URL")
	if itemsURL == "" {
		itemsURL = "http://localhost:8081"
	}

	return &service{
		repo:            repo,
		itemsServiceURL: itemsURL,
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (s *service) CreateOrder(req CreateOrderRequest) (*Order, error) {
	slog.Info("Validating Item in items-service", "item_id", req.ItemID, "quantity", req.Quantity)
	itemDetail, err := s.checkItemAvailability(req.ItemID)
	if err != nil {
		slog.Error("Item validation failed", "error", err)
		return nil, fmt.Errorf("%w: %v", ErrItemValidationFailed, err)
	}

	if itemDetail.Stock < req.Quantity {
		return nil, fmt.Errorf("%w. Available: %d, Requested: %d", ErrInsufficientStock, itemDetail.Stock, req.Quantity)
	}

	ord := Order{
		ID:       GenerateUUID(),
		UserID:   req.UserID,
		ItemID:   req.ItemID,
		Quantity: req.Quantity,
		Amount:   req.Amount,
		Address:  req.Address,
		Status:   StatusReadyToProcess,
	}

	if err := s.repo.Create(&ord); err != nil {
		return nil, fmt.Errorf("failed to persist order: %w", err)
	}

	slog.Info("Order successfully created", "order_id", ord.ID, "status", ord.Status)
	return &ord, nil
}

func (s *service) GetByID(id string) (*Order, error) {
	return s.repo.FindByID(id)
}

func (s *service) UpdateStatus(id, status string) error {
	if status != StatusPending && status != StatusReadyToProcess && status != StatusCompleted && status != StatusFailed {
		return ErrInvalidOrderStatus
	}

	return s.repo.UpdateStatus(id, status)
}

func (s *service) checkItemAvailability(itemID string) (*ItemResponseData, error) {
	url := fmt.Sprintf("%s/api/items/%s", s.itemsServiceURL, itemID)
	resp, err := s.httpClient.Get(url)
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
