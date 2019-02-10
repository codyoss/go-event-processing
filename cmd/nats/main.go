package main

import (
	"encoding/json"
	"fmt"

	"github.com/codyoss/go-event-processing/pkg/server"

	"github.com/nats-io/go-nats-streaming"
	"github.com/nats-io/go-nats-streaming/pb"

	"github.com/codyoss/go-event-processing/pkg/config"
)

func main() {
	c := config.Config{}
	config.Parse(&c)

	conn, _ := stan.Connect(c.Nats.ClusterID, "client-123")
	h := handlerProvider(c.Nats.OutputChannel, conn)
	conn.QueueSubscribe(c.Nats.InputChannel, c.Nats.Group, h, stan.DurableName("group"), stan.StartAt(pb.StartPosition_First))

	<-server.InterruptListener()
}

func handler(msg *stan.Msg) {
	fmt.Println(string(msg.Data))
}

func handlerProvider(outChannel string, conn stan.Conn) stan.MsgHandler {
	return func(msg *stan.Msg) {
		a := &Animal{}
		err := json.Unmarshal(msg.Data, a)
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
		fmt.Println(string(out))
		err = conn.Publish(outChannel, out)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// Animal is a creature on Earth
type Animal struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Age  int    `json:"age"`
}
