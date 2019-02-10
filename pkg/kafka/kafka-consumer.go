package kafka

import (
	"context"
	"time"

	kgo "github.com/segmentio/kafka-go"
)

// ConsumerConfig represents the config for a kafka consumer.
type ConsumerConfig struct {
	Brokers []string
	GroupID string
	Topic   string
}

// Consumer consumers messages from kafka.
type Consumer struct {
	rc     kgo.ReaderConfig
	reader *kgo.Reader
}

// NewConsumer creates a new Consumer.
func NewConsumer(c ConsumerConfig) *Consumer {
	return &Consumer{
		rc: kgo.ReaderConfig{
			Brokers:  c.Brokers,
			GroupID:  c.GroupID,
			Topic:    c.Topic,
			MinBytes: 1e3, // 1KB
			MaxBytes: 1e6, // 1MB
		},
	}
}

// Open up connections for reading from kafka.
func (c *Consumer) Open() {
	if c.reader == nil {
		c.reader = kgo.NewReader(c.rc)
	}
}

// Consume a message from kafka.
func (c *Consumer) Consume() ([]byte, error) {
	if c.reader == nil {
		panic("You must open the Consumer before you can receive messages.")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	m, err := c.reader.ReadMessage(ctx)
	return m.Value, err
}

// Close closes the kafka consumer.
func (c *Consumer) Close() {
	_ = c.reader.Close()
}
