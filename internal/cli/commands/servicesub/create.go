package servicesub

import (
	"NetManager/internal/minecraft"
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"strings"
)

type CreateSubcommand struct {
	Printer        interfaces.Printer
	ServiceManager interfaces.ServiceManager
}

func (cmd *CreateSubcommand) Execute(args []string) {
	if len(args) <= 0 {
		cmd.Printer.PrintColored("Service name not specified. 'service create <name>'", cmd.Printer.Service(), types.Yellow)
		return
	} else if len(args) > 1 {
		cmd.Printer.PrintColored("Service name cannot contain spaces.", cmd.Printer.Service(), types.Yellow)
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

func (cmd *CreateSubcommand) Subcommands() map[string]interfaces.Command {
	return map[string]interfaces.Command{}
}
