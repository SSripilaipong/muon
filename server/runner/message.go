package runner

import (
	"github.com/SSripilaipong/go-common/rslt"
	es "github.com/SSripilaipong/muon/server/eventsource"
)

type localAppendRequest struct {
	actions []es.Action
	reply   chan<- rslt.Of[es.AppendResponse]
}

func (r localAppendRequest) Actions() []es.Action                     { return r.actions }
func (r localAppendRequest) Reply() chan<- rslt.Of[es.AppendResponse] { return r.reply }

type markCommitUntilRequest struct {
	sequence uint64
	reply    chan<- error
}

func (r markCommitUntilRequest) Sequence() uint64    { return r.sequence }
func (r markCommitUntilRequest) Reply() chan<- error { return r.reply }
