package actor

import (
	"context"
	"log"

	"github.com/SSripilaipong/go-common/rslt"
)

type Processor[T any] interface {
	Process(msg T) rslt.Of[Processor[T]]
}

func StartLoop[T any](ctx context.Context, msgBox <-chan T, p Processor[T]) <-chan struct{} {
	done := make(chan struct{})

	var err error
	go func() {
		defer close(done)
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgBox:
				if !ok {
					continue
				}
				p, err = p.Process(msg).Return()
				if err != nil {
					log.Printf("actor processed message and got error: %s\n", err)
					return
				}
			}
		}
	}()
	return done
}
