package gateway

import (
	"github.com/SSripilaipong/muto/syntaxtree/result"

	"github.com/SSripilaipong/muon/server/runner"
)

type Runner interface { // TODO no need?
	Run(node result.SimplifiedNode) error
}

var _ Runner = runner.Service{}
