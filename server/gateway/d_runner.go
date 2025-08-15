package gateway

import (
	"context"

	"github.com/SSripilaipong/muto/syntaxtree/result"

	"github.com/SSripilaipong/muon/server/runner"
)

type Runner interface {
	Run(ctx context.Context, node result.SimplifiedNode) error
}

var _ Runner = runner.Service{}
