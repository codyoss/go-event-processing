package nats

import (
	"fmt"

	"github.com/nats-io/go-nats-streaming"
)

// ConsumerConfig provides configuration for a NATS consumer.
type ConsumerConfig struct {
	Channel   string
	ClusterID string
	Group     string
}

// Consumer consumers messages from NATS.
type Consumer struct {
	ch     chan []byte
	conn   stan.Conn
	config *ConsumerConfig
	sub    stan.Subscription
}

// NewConsumer creates a new Consumer.
func NewConsumer(c ConsumerConfig) *Consumer {
	conn, err := stan.Connect(c.ClusterID, "client-1")
	if err != nil {
		panic(fmt.Sprintf("could not connect to nats %s", err.Error()))
	}
	return &Consumer{
		ch:     make(chan []byte, 1),
		conn:   conn,
		config: &c,
	}
}

// Open a connection to NATS server.
func (c *Consumer) Open() {
	handler := func(msg *stan.Msg) {
		c.ch <- msg.Data
	}
	c.sub, _ = c.conn.QueueSubscribe(c.config.Channel, c.config.Group, handler, stan.DurableName(c.config.Group), stan.StartAt(0))
}

// Consume reads one msg from a NATS server.
func (c *Consumer) Consume() ([]byte, error) {
	return <-c.ch, nil
}

// Close closes the connection to NATS.
func (c *Consumer) Close() {
	_ = c.sub.Close()
	_ = c.conn.Close()
	close(c.ch)
}
