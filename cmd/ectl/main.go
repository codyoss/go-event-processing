package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/nats-io/go-nats-streaming"

	"github.com/codyoss/go-event-processing/pkg/kafka"
)

const (
	errKafkaOrNats        = "kafka or nats subcommand is required"
	errConsumerOrProducer = "consumer or producer subcommand is required"
)

var (
	broker string
	topic  string
	group  string
	msg    string

	cluster string
	channel string
)

func main() {
	kafkaConsumerCmd := flag.NewFlagSet("kafka-consumer", flag.ExitOnError)
	kafkaConsumerCmd.StringVar(&broker, "broker", "", "broker to consume from")
	kafkaConsumerCmd.StringVar(&topic, "topic", "", "topic to consume from")
	kafkaConsumerCmd.StringVar(&group, "group", "", "group id to consume from")

	kafkaProducerCmd := flag.NewFlagSet("kafka-producer", flag.ExitOnError)
	kafkaProducerCmd.StringVar(&broker, "broker", "", "broker to produce to")
	kafkaProducerCmd.StringVar(&topic, "topic", "", "topic to produce to")
	kafkaProducerCmd.StringVar(&msg, "msg", "", "a message to produce")

	natsConsumerCmd := flag.NewFlagSet("nats-consumer", flag.ExitOnError)
	natsConsumerCmd.StringVar(&cluster, "cluster", "", "cluster to produce to")
	natsConsumerCmd.StringVar(&channel, "channel", "", "channel to produce to")
	natsConsumerCmd.StringVar(&group, "group", "", "group id to consume from")
	natsProducerCmd := flag.NewFlagSet("nats-producer", flag.ExitOnError)
	natsProducerCmd.StringVar(&cluster, "cluster", "", "cluster to produce to")
	natsProducerCmd.StringVar(&channel, "channel", "", "channel to produce to")
	natsProducerCmd.StringVar(&msg, "msg", "", "a message to produce")

	if len(os.Args) < 2 {
		exit(errKafkaOrNats)
	}

	switch os.Args[1] {
	case "kafka":
		if isConsumer() {
			kafkaConsumerCmd.Parse(os.Args[3:])
		} else {
			kafkaProducerCmd.Parse(os.Args[3:])
		}
	case "nats":
		if isConsumer() {
			natsConsumerCmd.Parse(os.Args[3:])
		} else {
			natsProducerCmd.Parse(os.Args[3:])
		}
	default:
		exit(errKafkaOrNats)
	}

	switch {
	case kafkaConsumerCmd.Parsed():
		if broker == "" || topic == "" || group == "" {
			exit("broker, topic, and group must all be set")
		}
		consumeKafkaMsg()
	case kafkaProducerCmd.Parsed():
		if broker == "" || topic == "" || msg == "" {
			exit("broker, topic, and msg must all be set")
		}
		produceKafkaMsg()
	case natsConsumerCmd.Parsed():
		if channel == "" || cluster == "" || group == "" {
			exit("channel, cluster, and group must all be set")
		}
		consumeNatsMsg()
	case natsProducerCmd.Parsed():
		if channel == "" || cluster == "" || msg == "" {
			exit("channel, cluster, and msg must all be set")
		}
		produceNatsMsg()
	}
}

func isConsumer() bool {
	if len(os.Args) < 3 {
		exit(errConsumerOrProducer)
		os.Exit(1)
	}

	switch os.Args[2] {
	case "consumer":
		return true
	case "producer":
		return false
	default:
		exit(errConsumerOrProducer)
		return false
	}
}

func consumeKafkaMsg() {
	c := kafka.NewConsumer(kafka.ConsumerConfig{
		Brokers: []string{broker},
		GroupID: group,
		Topic:   topic,
	})
	c.Open()
	defer c.Close()
	for {
		b, err := c.Consume()
		if err != nil {
			return
		}
		fmt.Println(string(b))
	}
}
func produceKafkaMsg() {
	p := kafka.NewProducer(kafka.ProducerConfig{
		Brokers: []string{broker},
		Topic:   topic,
	})
	p.Open()
	defer p.Close()
	p.Produce([]byte(msg))
}

func consumeNatsMsg() {
	conn, err := stan.Connect(cluster, "client-198")
	if err != nil {
		exit(err.Error())
	}
	sub, err := conn.QueueSubscribe(channel, group, natsHandler, stan.DurableName(group), stan.StartAt(0))
	if err != nil {
		exit(err.Error())
	}
	time.Sleep(1000 * time.Millisecond)
	sub.Close()
	conn.Close()
}

func natsHandler(msg *stan.Msg) {
	fmt.Println(string(msg.Data))
}

func produceNatsMsg() {
	conn, err := stan.Connect(cluster, "client-198")
	if err != nil {
		exit(err.Error())
	}
	err = conn.Publish(channel, []byte(msg))
	if err != nil {
		exit(err.Error())
	}
	conn.Close()
}

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}
