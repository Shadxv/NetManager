package servicesub

import (
	"NetManager/external/cli"
	"NetManager/external/service"
	"NetManager/internal/minecraft"
	"strings"
)

type CreateSubcommand struct {
	Printer        cli.Printer
	ServiceManager service.Manager
}

func (cmd *CreateSubcommand) Execute(args []string) {
	if len(args) <= 0 {
		cmd.Printer.PrintColored("Service name not specified. 'service create <name>'", cmd.Printer.Service(), cli.Yellow)
		return
	} else if len(args) > 1 {
		cmd.Printer.PrintColored("Service name cannot contain spaces.", cmd.Printer.Service(), cli.Yellow)
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

func (cmd *CreateSubcommand) Subcommands() map[string]cli.Command {
	return map[string]cli.Command{}
}
