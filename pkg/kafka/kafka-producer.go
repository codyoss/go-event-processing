package kafka

import (
	"context"
	"time"

	kgo "github.com/segmentio/kafka-go"
)

// ProducerConfig is the configuration for a Producer.
type ProducerConfig struct {
	Brokers []string
	Topic   string
}

// Producer produces messages to kafka.
type Producer struct {
	wc     kgo.WriterConfig
	writer *kgo.Writer
}

// NewProducer creates a new Producer.
func NewProducer(c ProducerConfig) *Producer {
	return &Producer{
		wc: kgo.WriterConfig{
			Brokers:  c.Brokers,
			Topic:    c.Topic,
			Balancer: &kgo.RoundRobin{},
		},
	}
}

// Open opens a connection to kafka.
func (p *Producer) Open() {
	if p.writer == nil {
		p.writer = kgo.NewWriter(p.wc)
	}
}

// Produce produces a message to kafka.
func (p *Producer) Produce(msg []byte) error {
	if p.writer == nil {
		panic("You must open the Producer before you can send messages.")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	return p.writer.WriteMessages(ctx,
		kgo.Message{
			Value: msg,
		},
	)
}

// Close closes the connection to kafka.
func (p *Producer) Close() {
	p.writer.Close()
}
