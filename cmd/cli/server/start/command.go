package start

import (
	"github.com/urfave/cli/v2"

	"github.com/SSripilaipong/muon/server"
)

func NewCommand() *cli.Command {
	return &cli.Command{
		Name: "start",
		Action: func(ctx *cli.Context) error {
			return server.Start()
		},
	}
}
