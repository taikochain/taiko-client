package main

import (
	"os"

	"github.com/ethereum/go-ethereum/log"
	"github.com/taikochain/taiko-client/cmd/flags"
	"github.com/taikochain/taiko-client/cmd/utils"
	"github.com/taikochain/taiko-client/driver"
	"github.com/taikochain/taiko-client/proposer"
	"github.com/taikochain/taiko-client/prover"
	"github.com/urfave/cli/v2"
)

func main() {
	app := cli.NewApp()

	app.Version = "0.0.1"
	app.Name = "Taiko Clients"
	app.Description = "Entrypoint of Taiko Clients"
	app.Authors = []*cli.Author{{Name: "Taiko Labs", Email: "info@taiko.xyz"}}

	app.Commands = []*cli.Command{
		{
			Name:        "driver",
			Flags:       flags.DriverFlags,
			Description: "Taiko Driver software",
			Action:      utils.SubcommandAction(new(driver.Driver)),
		},
		{
			Name:        "proposer",
			Flags:       flags.ProposerFlags,
			Description: "Taiko Proposer software",
			Action:      utils.SubcommandAction(new(proposer.Proposer)),
		},
		{
			Name:        "prover",
			Flags:       flags.ProverFlags,
			Description: "Taiko Prover software",
			Action:      utils.SubcommandAction(new(prover.Prover)),
		},
	}

	app.Action = func(ctx *cli.Context) error {
		log.Crit("Expected driver/proposer/prover subcommands")
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Crit("Failed to start Taiko client", "error", err)
	}
}
