package rabbitmq

import (
	"context"
	"fmt"
	"github/socialforge/config"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// Queue names for the ingestion pipeline.
const (
	QueueIngestInbound    = "ingest.inbound"
	QueueDispatchOutbound = "dispatch.outbound"
)

// Client is a thin RabbitMQ wrapper: one connection, a dedicated publish
// channel, and one channel per consumer. It declares durable queues on connect.
type Client struct {
	conn   *amqp.Connection
	pubCh  *amqp.Channel
	cfg    *config.RabbitMQConfig
	logger *zap.Logger
	mu     sync.Mutex
}

func NewClient(ctx context.Context, cfg *config.RabbitMQConfig, logger *zap.Logger) (*Client, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.VHost)

	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %w", err)
	}
	pubCh, err := conn.Channel()
	if err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("failed to open publish channel: %w", err)
	}

	c := &Client{conn: conn, pubCh: pubCh, cfg: cfg, logger: logger}
	if err := c.declareTopology(); err != nil {
		_ = c.Close()
		return nil, err
	}

	logger.Info("✅ RabbitMQ client initialized successfully", zap.String("host", cfg.Host), zap.Int("port", cfg.Port))
	return c, nil
}

func (c *Client) declareTopology() error {
	for _, q := range []string{QueueIngestInbound, QueueDispatchOutbound} {
		if _, err := c.pubCh.QueueDeclare(q, true, false, false, false, nil); err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", q, err)
		}
	}
	return nil
}

// Publish sends a persistent message to a queue (via the default exchange).
func (c *Client) Publish(ctx context.Context, queue string, body []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.pubCh.PublishWithContext(ctx, "", queue, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Body:         body,
	})
}

// Consume starts a background consumer for a queue. A handler error causes one
// requeue; a message that fails again (redelivered) is dropped to avoid poison
// loops (ingestion is idempotent via dedup, so a requeue is safe).
func (c *Client) Consume(queue, consumerTag string, prefetch int, handler func(ctx context.Context, body []byte) error) error {
	ch, err := c.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open consumer channel: %w", err)
	}
	if _, err := ch.QueueDeclare(queue, true, false, false, false, nil); err != nil {
		_ = ch.Close()
		return fmt.Errorf("failed to declare queue %s: %w", queue, err)
	}
	if prefetch <= 0 {
		prefetch = 10
	}
	if err := ch.Qos(prefetch, 0, false); err != nil {
		_ = ch.Close()
		return fmt.Errorf("failed to set qos: %w", err)
	}

	deliveries, err := ch.Consume(queue, consumerTag, false, false, false, false, nil)
	if err != nil {
		_ = ch.Close()
		return fmt.Errorf("failed to start consumer: %w", err)
	}

	go func() {
		defer ch.Close()
		for d := range deliveries {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			err := handler(ctx, d.Body)
			cancel()

			if err == nil {
				_ = d.Ack(false)
				continue
			}
			c.logger.Error("consumer handler failed",
				zap.String("queue", queue), zap.Bool("redelivered", d.Redelivered), zap.Error(err))
			// Requeue once; drop on a second failure to avoid poison loops.
			_ = d.Nack(false, !d.Redelivered)
		}
		c.logger.Warn("rabbitmq consumer channel closed", zap.String("queue", queue))
	}()

	c.logger.Info("🐇 consumer started", zap.String("queue", queue), zap.String("tag", consumerTag))
	return nil
}

func (c *Client) Close() error {
	if c.pubCh != nil {
		_ = c.pubCh.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}
