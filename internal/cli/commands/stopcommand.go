package commands

import (
	"NetManager/internal/cli/model"
)

type StopCommand struct {
	Printer model.Printer
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
