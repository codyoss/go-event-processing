package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/lovoo/goka/codec"

	"github.com/gorilla/mux"
	"github.com/lovoo/goka"
)

var (
	brokers                 = []string{"localhost:9092"}
	inputTopic  goka.Stream = "output-topic"
	outputTopic goka.Stream = "goka-output"
	group       goka.Group  = "group"
)

func main() {
	go runProcessor()
	go runProcessor2()
	runView()
}

func runProcessor() {
	g := goka.DefineGroup(group,
		goka.Input(inputTopic, new(animalCodec), processOne),
		goka.Output(outputTopic, new(codec.Int64)),
	)
	p, err := goka.NewProcessor(brokers, g)
	if err != nil {
		panic(err)
	}

	p.Run(context.Background())
}
func runProcessor2() {
	g := goka.DefineGroup(group,
		goka.Input(outputTopic, new(codec.Int64), processTwo),
		goka.Persist(new(animalCodec)),
	)
	p, err := goka.NewProcessor(brokers, g)
	if err != nil {
		panic(err)
	}

	p.Run(context.Background())
}

func runView() {
	view, err := goka.NewView(brokers,
		goka.GroupTable(group),
		new(animalCodec),
	)
	if err != nil {
		panic(err)
	}

	root := mux.NewRouter()
	root.HandleFunc("/{key}", func(w http.ResponseWriter, r *http.Request) {
		value, _ := view.Get(mux.Vars(r)["key"])
		data, _ := json.Marshal(value)
		w.Write(data)
	})
	fmt.Println("View opened at http://localhost:9095/")
	go http.ListenAndServe(":9095", root)

	view.Run(context.Background())
}

func processOne(ctx goka.Context, msg interface{}) {
	a := msg.(*Animal)
	ctx.Emit(outputTopic, a.Type, a.Age)
	fmt.Printf("name: %s age: %d\n", a.Name, a.Age)
}

func processTwo(ctx goka.Context, msg interface{}) {
	var a *Animal
	if val := ctx.Value(); val != nil {
		a = val.(*Animal)
	} else {
		a = new(Animal)
		a.Type = ctx.Key()
	}
	age := msg.(int64)

	a.Age += age
	ctx.SetValue(a)
	fmt.Printf("key: %s total age: %d\n", ctx.Key(), a.Age)
}

// Animal is a creature on Earth
type Animal struct {
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Age  int64  `json:"age"`
}

type animalCodec struct{}

// Encodes an Animal into []byte
func (jc *animalCodec) Encode(value interface{}) ([]byte, error) {
	if _, ok := value.(*Animal); !ok {
		return nil, fmt.Errorf("Codec requires value *Animal, got %T", value)
	}
	return json.Marshal(value)
}

// Decodes an Animal from []byte to it's go representation.
func (jc *animalCodec) Decode(data []byte) (interface{}, error) {
	var (
		a   Animal
		err error
	)
	err = json.Unmarshal(data, &a)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling Animal: %v", err)
	}
	return &a, nil
}
