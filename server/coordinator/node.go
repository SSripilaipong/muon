package coordinator

import (
	"context"

	es "github.com/SSripilaipong/muon/server/eventsource"
)

type Node interface {
	Commit(ctx context.Context, actions []es.Action) error
}
