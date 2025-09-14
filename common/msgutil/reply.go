package msgutil

import (
	"time"

	"github.com/SSripilaipong/muon/common/chn"
)

type ReplyMixin[T any] struct {
	ch      chan T
	timeout time.Duration
}

func NewReplyMixin[T any](ch chan T, timeout time.Duration) ReplyMixin[T] {
	return ReplyMixin[T]{ch: ch, timeout: timeout}
}

func (r ReplyMixin[T]) Reply(x T) error {
	return chn.SendWithTimeout(r.ch, x, r.timeout)
}
