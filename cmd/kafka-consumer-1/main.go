package main

import (
	"context"
	"fmt"
	"time"

	"github.com/codyoss/go-event-processing/pkg/config"

	"github.com/segmentio/kafka-go"
)

func main() {
	c := &config.Config{}
	config.Parse(c)

	rc := kafka.ReaderConfig{
		Brokers: c.Kafka.Brokers,
		GroupID: c.Kafka.Group,
		Topic:   c.Kafka.InputTopic,
	}

	r := kafka.NewReader(rc)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	m, _ := r.ReadMessage(ctx)
	fmt.Println(string(m.Value))
	r.Close()
}
