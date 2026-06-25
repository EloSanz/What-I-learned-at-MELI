package broker

import (
	"encoding/json"
	"log/slog"

	sdkbroker "github.com/user/meli-sdk/broker"
	"orders-service/internal/order"
)

type Consumer struct {
	rabbitMQ     *sdkbroker.RabbitMQ
	orderService order.Service
}

func NewConsumer(rabbitMQ *sdkbroker.RabbitMQ, orderService order.Service) *Consumer {
	return &Consumer{
		rabbitMQ:     rabbitMQ,
		orderService: orderService,
	}
}

func (c *Consumer) Start() {
	// Listen for stock.reserved
	err := c.rabbitMQ.Consume("orders_stock_reserved_queue", "stock.reserved", func(body []byte) error {
		var event map[string]interface{}
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		orderID, _ := event["order_id"].(string)
		if orderID != "" {
			slog.Info("Received stock.reserved event, updating order status", "order_id", orderID)
			return c.orderService.UpdateStatus(orderID, order.StatusReadyToProcess)
		}
		return nil
	})
	if err != nil {
		slog.Error("Failed to start stock.reserved consumer", "error", err)
	}

	// Listen for stock.rejected
	err = c.rabbitMQ.Consume("orders_stock_rejected_queue", "stock.rejected", func(body []byte) error {
		var event map[string]interface{}
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}
		orderID, _ := event["order_id"].(string)
		if orderID != "" {
			slog.Info("Received stock.rejected event, updating order status", "order_id", orderID)
			return c.orderService.UpdateStatus(orderID, order.StatusFailed)
		}
		return nil
	})
	if err != nil {
		slog.Error("Failed to start stock.rejected consumer", "error", err)
	}
}
