package servicesub

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"strings"
)

type BuildSubcommand struct {
	Printer        interfaces.Printer
	ServiceManager interfaces.ServiceManager
	ImageManager   interfaces.ImageManager
}

func (cmd *BuildSubcommand) Execute(args []string) {
	if len(args) <= 0 {
		cmd.Printer.PrintColored("Service name not specified. 'service create <name>'", cmd.Printer.Service(), types.Yellow)
		return
	} else if len(args) > 1 {
		cmd.Printer.PrintColored("Service name cannot contain spaces.", cmd.Printer.Service(), types.Yellow)
		return
	}

	name := strings.ToLower(args[0])
	serviceModel := cmd.ServiceManager.GetService(name)
	if serviceModel == nil {
		cmd.Printer.PrintColored("Service with this name does not exist.", cmd.Printer.Service(), types.Yellow)
		return
	}

	serviceModel.Build(cmd.Printer, cmd.ImageManager, cmd.ServiceManager)
}

func (cmd *BuildSubcommand) Name() string {
	return "build"
}

func (cmd *BuildSubcommand) Description() string {
	return "Builds image for specific service"
}

func (cmd *BuildSubcommand) Subcommands() map[string]interfaces.Command {
	return map[string]interfaces.Command{}
}
