package servicesub

import (
	"NetManager/external/cli"
	"NetManager/external/docker"
	"NetManager/external/service"
	"strings"
)

type BuildSubcommand struct {
	Printer        cli.Printer
	ServiceManager service.Manager
	ImageManager   docker.ImageManager
}

func (cmd *BuildSubcommand) Execute(args []string) {
	if len(args) <= 0 {
		cmd.Printer.PrintColored("Service name not specified. 'service create <name>'", cmd.Printer.Service(), cli.Yellow)
		return
	} else if len(args) > 1 {
		cmd.Printer.PrintColored("Service name cannot contain spaces.", cmd.Printer.Service(), cli.Yellow)
		return
	}

	name := strings.ToLower(args[0])
	serviceModel := cmd.ServiceManager.GetService(name)
	if serviceModel == nil {
		cmd.Printer.PrintColored("Service with this name does not exist.", cmd.Printer.Service(), cli.Yellow)
		return
	}

	cmd.ImageManager.FullDeployImage(serviceModel, cmd.ServiceManager)
}

func (cmd *BuildSubcommand) Name() string {
	return "build"
}

func (cmd *BuildSubcommand) Description() string {
	return "Builds image for specific service"
}

func (cmd *BuildSubcommand) Subcommands() map[string]cli.Command {
	return map[string]cli.Command{}
}
