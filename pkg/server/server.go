package server

import (
	"os"
	"os/signal"
	"sync"
	"time"
)

// InterruptListener returns a blocking channel. It will unblock after sigInt.
func InterruptListener() <-chan os.Signal {
	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, os.Interrupt)
	return sigInt
}

// Wait will wait for WaitGroup to unblock. This method will return after timeout is exceeded if the WaitGroup is not
// done waiting. Returns true if the waiting on the WaitGroup was successful.
func Wait(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		wg.Wait()
		close(c)
	}()
	select {
	case <-c:
		return true
	case <-time.After(timeout):
		return false
	}
}
