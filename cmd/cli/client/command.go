package client

import (
	"github.com/urfave/cli/v2"

	"github.com/SSripilaipong/muon/cmd/cli/client/repl"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name: "client",
		Subcommands: []*cli.Command{
			repl.NewCommand(),
		},
	}
}
