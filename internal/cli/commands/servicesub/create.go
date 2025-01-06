package servicesub

import (
	"NetManager/internal/cli/model"
	"NetManager/internal/minecraft"
	types "NetManager/internal/minecraft/type"
	"strings"
)

type CreateSubcommand struct {
	Printer        model.Printer
	ServiceManager types.ServiceManagerBase
}

func (cmd *CreateSubcommand) Execute(args []string) {
	if len(args) <= 0 {
		cmd.Printer.PrintColored("Service name not specified. 'service create <name>'", cmd.Printer.Service(), model.Yellow)
		return
	} else if len(args) > 1 {
		cmd.Printer.PrintColored("Service name cannot contain spaces.", cmd.Printer.Service(), model.Yellow)
		return
	}

	name := strings.ToLower(args[0])
	minecraft.NetServiceWizard(cmd.Printer, cmd.ServiceManager, name).Run()
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
