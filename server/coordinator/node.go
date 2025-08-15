package coordinator

import (
	es "github.com/SSripilaipong/muon/server/eventsource"
)

type Node interface {
	Commit(actions []es.Action) error
}
