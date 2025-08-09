package main

import (
	"os"

	"github.com/urfave/cli/v2"

	"github.com/SSripilaipong/muon/cmd/cli/client"
	"github.com/SSripilaipong/muon/cmd/cli/server"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			server.NewCommand(),
			client.NewCommand(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
