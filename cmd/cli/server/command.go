package server

import (
	"github.com/urfave/cli/v2"

	"github.com/SSripilaipong/muon/cmd/cli/server/start"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name: "server",
		Subcommands: []*cli.Command{
			start.NewCommand(),
		},
	}
}
