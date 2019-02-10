package message

import "fmt"

// Consumer receives messages from a messaging system.
type Consumer interface {
	Open()
	Consume() ([]byte, error)
	Close()
}

// Read keeps on receiving messages until the done channel is closed. This method will open and close the consumer for
// you. The supplied function will operate on the bytes received from the Consumer. If there is an error of any kind it
// will be passed to the error in fn. If a fn is not provided the default behavior is to print the byte[] as a string
// or the error if there is one.
func Read(c Consumer, done <-chan struct{}, fn func([]byte, error)) {
	if fn == nil {
		fn = defaultFn
	}
	c.Open()
	read(c, done, fn)
	c.Close()
}

func read(c Consumer, done <-chan struct{}, fn func([]byte, error)) {
	for {
		select {
		default:
			fn(c.Consume())
		case <-done:
			return
		}
	}
}

func defaultFn(b []byte, err error) {
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	}
	fmt.Println(string(b))
}
