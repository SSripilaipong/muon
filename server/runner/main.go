package runner

import (
	"github.com/SSripilaipong/muto/syntaxtree/result"
)

type Runner struct{}

func New() Runner {
	return Runner{}
}

func (Runner) Start() error {
	return nil
}

func (Runner) Stop() error {
	return nil
}

func (p Runner) Done() chan struct{} {
	return nil
}

func (p Runner) Run(node result.Node) error {
	return nil
}
