package nats

import (
	"fmt"

	"github.com/nats-io/go-nats-streaming"
)

// ProducerConfig is the configuration for a Producer.
type ProducerConfig struct {
	Channel   string
	ClusterID string
}

// Producer produces messages to NATS.
type Producer struct {
	conn   stan.Conn
	config *ProducerConfig
}

// NewProducer creates a new NATS producer.
func NewProducer(c ProducerConfig) *Producer {
	return &Producer{
		config: &c,
	}
}

// Open a connection to NATS server.
func (p *Producer) Open() {
	conn, err := stan.Connect(p.config.ClusterID, "producer-1")
	if err != nil {
		panic(fmt.Sprintf("could not connect to nats %s", err.Error()))
	}
	p.conn = conn
}

// Produce message to NATS server.
func (p *Producer) Produce(msg []byte) error {
	return p.conn.Publish(p.config.Channel, msg)
}

// Close connection to NATS server.
func (p *Producer) Close() {
	p.conn.Close()
}
