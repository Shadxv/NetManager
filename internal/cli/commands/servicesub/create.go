package servicesub

import "NetManager/internal/cli/model"

type CreateSubcommand struct {
	Printer model.Printer
}

func (cmd *CreateSubcommand) Execute(args []string) {
	cmd.Printer.Print("Create", cmd.Printer.Service())
}

func (cmd *CreateSubcommand) Name() string {
	return "create"
}

func (cmd *CreateSubcommand) Description() string {
	return "Creates new service"
}

func (cmd *CreateSubcommand) Subcommands() map[string]model.Command {
	return map[string]model.Command{}
}
