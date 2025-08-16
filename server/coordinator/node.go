package coordinator

import (
	"context"

	es "github.com/SSripilaipong/muon/server/eventsource"
)

type Node interface {
	Append(ctx context.Context, actions []es.Action) error
}

type LocalNode interface {
	Node
	LocalAppend(ctx context.Context, actions []es.Action) error
}
