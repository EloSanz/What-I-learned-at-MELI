package order

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/user/meli-sdk/broker"
	"github.com/user/meli-sdk/lock"
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
	repo        Repository
	lockService lock.Service
	rabbitMQ    *broker.RabbitMQ
}

func NewService(repo Repository, lockService lock.Service, rabbitMQ *broker.RabbitMQ) Service {
	return &service{
		repo:        repo,
		lockService: lockService,
		rabbitMQ:    rabbitMQ,
	}
}

func (s *service) CreateOrder(req CreateOrderRequest) (*Order, error) {

	ord := Order{
		ID:       GenerateUUID(),
		UserID:   req.UserID,
		ItemID:   req.ItemID,
		Quantity: req.Quantity,
		Amount:   req.Amount,
		Address:  req.Address,
		Status:   StatusPending,
	}

	if err := s.repo.Create(&ord); err != nil {
		return nil, fmt.Errorf("failed to persist order: %w", err)
	}

	// Publish async event to RabbitMQ
	event := map[string]interface{}{
		"order_id": ord.ID,
		"item_id":  ord.ItemID,
		"quantity": ord.Quantity,
	}
	if err := s.rabbitMQ.Publish(context.Background(), "order.created", event); err != nil {
		slog.Error("Failed to publish order.created event", "order_id", ord.ID, "error", err)
	}

	slog.Info("Order successfully created and queued for processing", "order_id", ord.ID, "status", ord.Status)
	return &ord, nil
}

func (s *service) GetByID(id string) (*Order, error) {
	return s.repo.FindByID(id)
}

func (s *service) UpdateStatus(id, status string) error {
	if status != StatusPending && status != StatusReadyToProcess && status != StatusCompleted && status != StatusFailed {
		return ErrInvalidOrderStatus
	}

	return s.lockService.WithLock(id, func() error {
		return s.repo.UpdateStatus(id, status)
	})
}


