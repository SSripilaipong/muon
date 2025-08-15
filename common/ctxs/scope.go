package ctxs

import (
	"context"
	"time"
)

func TimeoutScope(ctx context.Context, timeout time.Duration, f func(context.Context)) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	f(ctxWithTimeout)
}
