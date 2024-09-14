package commands

import (
	"NetManager/internal/cli/model"
)

type StopCommand struct {
	Printer model.Printer
}

func (cmd *StopCommand) Execute(args []string) {
	cmd.Printer.SetClosingStatus()
	cmd.Printer.Print("Closing console...", cmd.Printer.Service())
	cmd.Printer.Close()
}

func (cmd *StopCommand) Name() string {
	return "stop"
}

func (cmd *StopCommand) Description() string {
	return "Stops app and all services"
}
