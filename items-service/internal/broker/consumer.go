package broker

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	sdkbroker "github.com/user/meli-sdk/broker"
	"github.com/user/meli-sdk/lock"
	"items-service/internal/item"
)

type Consumer struct {
	rabbitMQ    *sdkbroker.RabbitMQ
	itemService item.Service
}

func NewConsumer(rabbitMQ *sdkbroker.RabbitMQ, itemService item.Service) *Consumer {
	return &Consumer{
		rabbitMQ:    rabbitMQ,
		itemService: itemService,
	}
}

func (c *Consumer) Start() {
	err := c.rabbitMQ.Consume("items_order_created_queue", "order.created", func(body []byte) error {
		var event map[string]interface{}
		if err := json.Unmarshal(body, &event); err != nil {
			return err // Parse error, nack
		}

		orderID, _ := event["order_id"].(string)
		itemID, _ := event["item_id"].(string)
		quantityFloat, _ := event["quantity"].(float64)
		quantity := int(quantityFloat)

		if orderID == "" || itemID == "" || quantity <= 0 {
			slog.Warn("Invalid order.created event payload")
			return nil // Drop it (ack)
		}

		slog.Info("Processing order.created event", "order_id", orderID, "item_id", itemID)

		_, err := c.itemService.ValidateAndReserveStock(itemID, quantity)
		if err != nil {
			if errors.Is(err, lock.ErrResourceLocked) {
				slog.Warn("Resource locked, will retry later", "item_id", itemID)
				return err // Nack and requeue so RabbitMQ tries again
			}
			
			// It failed permanently (ErrOutOfStock, ErrItemNotFound)
			slog.Error("Stock reservation failed", "order_id", orderID, "error", err)
			
			// Publish rejection
			respEvent := map[string]interface{}{"order_id": orderID, "reason": err.Error()}
			c.rabbitMQ.Publish(context.Background(), "stock.rejected", respEvent)
			return nil // Ack, we handled the rejection
		}

		// Success!
		slog.Info("Stock reserved successfully for order", "order_id", orderID)
		respEvent := map[string]interface{}{"order_id": orderID}
		c.rabbitMQ.Publish(context.Background(), "stock.reserved", respEvent)
		return nil
	})

	if err != nil {
		slog.Error("Failed to start order.created consumer", "error", err)
	}
}
