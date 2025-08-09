package repl

import (
	"github.com/urfave/cli/v2"

	"github.com/SSripilaipong/muon/client/repl"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name:  "repl",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) error {
			return loop(repl.NewRepl())
		},
	}
}
