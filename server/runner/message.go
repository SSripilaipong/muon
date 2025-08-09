package runner

import (
	"github.com/SSripilaipong/muto/syntaxtree/result"
)

type runMessage struct {
	moduleVersion string
	node          result.SimplifiedNode
	reply         chan<- error
}

func (m runMessage) Reply() chan<- error   { return m.reply }
func (m runMessage) ModuleVersion() string { return m.moduleVersion }
