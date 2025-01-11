package commands

import (
	"NetManager/external/cli"
)

type StopCommand struct {
	Printer cli.Printer
}

func (cmd *StopCommand) Execute(args []string) {
	cmd.Printer.CloseGracefully("Closing console...")
}

func (cmd *StopCommand) Name() string {
	return "stop"
}

func (cmd *StopCommand) Description() string {
	return "Stops app and all services"
}

func (cmd *StopCommand) Subcommands() map[string]cli.Command {
	return map[string]cli.Command{}
}
