package repl

import (
	"github.com/SSripilaipong/go-common/optional"
)

type Repl[Cmd any] interface {
	Read() optional.Of[Cmd]
	Execute(Cmd) error
}

func loop[Cmd any](p Repl[Cmd]) error {
	for {
		cmd := p.Read()
		if cmd.IsEmpty() {
			continue
		}

		err := p.Execute(cmd.Value())

		if err != nil {
			return err
		}
	}
}
