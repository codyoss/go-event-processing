package main

import (
	"context"
	"encoding/json"
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
	consumer := createConsumer(c)
	producer := createProducer(c)
	producer.Open()
	processFn := processor(producer)
	for i := 0; i < c.App.Parallelism; i++ {
		wg.Add(1)
		go func() {
			message.Read(consumer, done, processFn)
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
	producer.Close()
	fmt.Println("Gracefully shutdown")
}

// converts human years into animal years
func processor(p message.Producer) func([]byte, error) {
	return func(b []byte, err error) {
		if err != nil {
			if err == context.DeadlineExceeded {
				return
			}
			fmt.Printf("Error: %s\n", err.Error())
		}
		a := &Animal{}
		err = json.Unmarshal(b, a)
		if err != nil {
			fmt.Println(err)
			return
		}
		switch {
		case a.Type == "Dog":
			a.Age = a.Age * 7
		case a.Type == "Cat":
			a.Age = a.Age * 6
		}

		out, _ := json.Marshal(a)
		err = p.Produce(out)
		if err != nil {
			fmt.Printf("Error: %s\n", err.Error())
			return
		}
		fmt.Println("Processed a msg")
	}
}

func createConsumer(c *config.Config) *kafka.Consumer {
	return kafka.NewConsumer(kafka.ConsumerConfig{
		Brokers: c.Kafka.Brokers,
		GroupID: c.Kafka.Group,
		Topic:   c.Kafka.InputTopic,
	})
}

func createProducer(c *config.Config) *kafka.Producer {
	return kafka.NewProducer(kafka.ProducerConfig{
		Brokers: c.Kafka.Brokers,
		Topic:   c.Kafka.OutputTopic,
	})
}

// Animal is a creature on Earth
type Animal struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Age  int    `json:"age"`
}
