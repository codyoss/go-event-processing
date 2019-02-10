package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/codyoss/go-event-processing/pkg/config"
	"github.com/codyoss/go-event-processing/pkg/kafka"
	"github.com/codyoss/go-event-processing/pkg/message"
	"github.com/codyoss/go-event-processing/pkg/server"
)

func main() {
	c := &config.Config{}
	config.Parse(c)

	done := make(chan struct{})
	wg := &sync.WaitGroup{}
	cc := consumerConfig(c)
	for i := 0; i < c.App.Parallelism; i++ {
		wg.Add(1)
		go func() {
			message.Read(kafka.NewConsumer(cc), done, nil)
			wg.Done()
		}()
	}

	<-server.InterruptListener()
	close(done)
	fmt.Printf("\nShutting down\n")

	if !server.Wait(wg, time.Duration(c.App.ShutdownTimeoutSec)*time.Second) {
		fmt.Println("Exceeded shutdown timeout")
		os.Exit(1)
	}
	fmt.Println("Gracefully shutdown")
}

func consumerConfig(c *config.Config) kafka.ConsumerConfig {
	return kafka.ConsumerConfig{
		Brokers: c.Kafka.Brokers,
		GroupID: c.Kafka.Group,
		Topic:   c.Kafka.InputTopic,
	}
}
