package repl

import "github.com/SSripilaipong/go-common/optional"

type Repl struct{}

func NewRepl() Repl {
	return Repl{}
}

func (r Repl) Read() optional.Of[string] {
	return optional.Empty[string]()
}

func (r Repl) Execute(string) error {
	return nil
}
