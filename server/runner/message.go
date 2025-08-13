package runner

import (
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
