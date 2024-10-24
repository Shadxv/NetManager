package commands

import (
	"NetManager/internal/cli/commands/servicesub"
	"NetManager/internal/cli/handler"
	"NetManager/internal/cli/model"
)

type ServiceCommand struct {
	Printer model.Printer
}

func (cmd *ServiceCommand) Execute(args []string) {
	handler.HandleSubcommand(cmd.Printer, cmd, args)
}

func (cmd *ServiceCommand) Name() string {
	return "service"
}

func (cmd *ServiceCommand) Description() string {
	return "Command to manage all services on NetManager"
}

func (cmd *ServiceCommand) Subcommands() map[string]model.Command {
	createSubcommand := servicesub.CreateSubcommand{
		Printer: cmd.Printer,
	}

	return map[string]model.Command{
		createSubcommand.Name(): &createSubcommand,
	}
}
