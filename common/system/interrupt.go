package system

import (
	"os"
	"os/signal"
)

func WaitForInterrupt() <-chan struct{} {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	out := make(chan struct{})
	go func() {
		<-c
		out <- struct{}{}
	}()
	return out
}
