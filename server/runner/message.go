package runner

import (
	"github.com/SSripilaipong/go-common/rslt"
	"github.com/SSripilaipong/muon/server/coordinator"
	es "github.com/SSripilaipong/muon/server/eventsource"
	"github.com/SSripilaipong/muto/syntaxtree/result"
)

type runRequest struct {
	moduleVersion string
	node          result.SimplifiedNode
	reply         chan<- error
}

func (r runRequest) ModuleVersion() string       { return r.moduleVersion }
func (r runRequest) Node() result.SimplifiedNode { return r.node }
func (r runRequest) Reply() chan<- error         { return r.reply }

type setCoordinatorRequest struct {
	coord *coordinator.Controller
}

func (r setCoordinatorRequest) Coordinator() *coordinator.Controller { return r.coord }

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
