package coordinator

import (
	"context"

	"github.com/SSripilaipong/go-common/rslt"

	es "github.com/SSripilaipong/muon/server/eventsource"
)

type Node interface {
	ForceAppend(ctx context.Context, actions []es.Action) rslt.Of[es.AppendResponse]
}

type LocalNode interface {
	LocalAppend(ctx context.Context, actions []es.Action) rslt.Of[es.AppendResponse]
	MarkCommitUntil(ctx context.Context, sequence uint64) error
}
