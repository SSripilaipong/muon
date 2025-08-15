package chn

import (
	"context"
	"fmt"
	"time"

	"github.com/SSripilaipong/go-common/optional"
	"github.com/SSripilaipong/go-common/rslt"
)

func ReceiveNoWait[T any](c <-chan T) optional.Of[T] {
	select {
	case x, ok := <-c:
		if ok {
			return optional.Value(x)
		}
	default:
	}
	return optional.Empty[T]()
}

func ReceiveWithTimeout[T any](c <-chan T, timeout time.Duration) rslt.Of[T] {
	select {
	case x, ok := <-c:
		if !ok {
			return rslt.Error[T](fmt.Errorf("cannot retrieve from channel"))
		}
		return rslt.Value(x)
	case <-time.After(timeout):
		return rslt.Error[T](fmt.Errorf("timed out"))
	}
}

func ReceiveWithContext[T any](ctx context.Context, c <-chan T) rslt.Of[T] {
	select {
	case x, ok := <-c:
		if !ok {
			return rslt.Error[T](fmt.Errorf("cannot retrieve from channel"))
		}
		return rslt.Value(x)
	case <-ctx.Done():
		return rslt.Error[T](ctx.Err())
	}
}
