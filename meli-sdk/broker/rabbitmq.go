package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

// NewRabbitMQ creates a new RabbitMQ connection and channel with retries
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	var conn *amqp.Connection
	var err error

	// Retry up to 10 times, 2 seconds apart
	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial(url)
		if err == nil {
			break
		}
		slog.Warn("Failed to connect to RabbitMQ, retrying...", "attempt", i+1, "error", err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ after retries: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	// Declare a standard topic exchange
	err = ch.ExchangeDeclare(
		"meli_events", // name
		"topic",       // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &RabbitMQ{conn: conn, ch: ch}, nil
}

// Publish JSON-encodes and sends a message to the given routing key
func (r *RabbitMQ) Publish(ctx context.Context, routingKey string, body interface{}) error {
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to encode body: %w", err)
	}

	err = r.ch.PublishWithContext(ctx,
		"meli_events", // exchange
		routingKey,    // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        b,
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	slog.Info("Message published", "routingKey", routingKey)
	return nil
}

// Consume starts a consumer on a queue bound to the given routing key
func (r *RabbitMQ) Consume(queueName string, routingKey string, handler func(body []byte) error) error {
	q, err := r.ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %w", err)
	}

	err = r.ch.QueueBind(
		q.Name,        // queue name
		routingKey,    // routing key
		"meli_events", // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	msgs, err := r.ch.Consume(
		q.Name, // queue
		"",     // consumer tag
		false,  // auto-ack (we use manual ack for reliability)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	go func() {
		for d := range msgs {
			err := handler(d.Body)
			if err != nil {
				slog.Error("Failed to process message", "routingKey", routingKey, "error", err)
				// Nack and requeue the message on failure
				d.Nack(false, true)
			} else {
				// Ack the message on success
				d.Ack(false)
			}
		}
	}()

	slog.Info("Started consumer", "queue", queueName, "routingKey", routingKey)
	return nil
}

// Close gracefully closes the channel and connection
func (r *RabbitMQ) Close() {
	if r.ch != nil {
		r.ch.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
