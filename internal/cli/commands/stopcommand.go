package commands

import (
	"NetManager/pkg/interfaces"
)

type StopCommand struct {
	Printer interfaces.Printer
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

func (cmd *StopCommand) Subcommands() map[string]interfaces.Command {
	return map[string]interfaces.Command{}
}
