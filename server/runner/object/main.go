package object

import (
	"context"

	"github.com/SSripilaipong/muto/core/base"

	"github.com/SSripilaipong/muon/common/chn"
)

func Spawn(ctx context.Context, node base.Node) {
	go func() {
		for base.IsMutableNode(node) && chn.ReceiveNoWait(ctx.Done()).IsEmpty() {
			result := base.UnsafeNodeToMutable(node).Mutate()
			if result.IsEmpty() {
				return
			}

			node = result.Value()
		}
	}()
}
