package commands

import (
	"NetManager/external/cli"
	"NetManager/external/docker"
	"NetManager/external/kubernetes"
	"NetManager/external/service"
	"NetManager/internal/cli/commands/servicesub"
	"NetManager/internal/cli/handler"
)

type ServiceCommand struct {
	Printer          cli.Printer
	ServiceManager   service.Manager
	ImageManager     docker.ImageManager
	KubernetesClient kubernetes.Client
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

func (cmd *ServiceCommand) Subcommands() map[string]cli.Command {
	createSubcommand := servicesub.CreateSubcommand{
		Printer:        cmd.Printer,
		ServiceManager: cmd.ServiceManager,
	}

	listSubcommand := servicesub.ListSubcommand{
		Printer:        cmd.Printer,
		ServiceManager: cmd.ServiceManager,
	}
	buildSubcommand := servicesub.BuildSubcommand{
		Printer:        cmd.Printer,
		ServiceManager: cmd.ServiceManager,
		ImageManager:   cmd.ImageManager,
	}
	startSubcommand := servicesub.StartSubcommand{
		Printer:          cmd.Printer,
		ServiceManager:   cmd.ServiceManager,
		KubernetesClient: cmd.KubernetesClient,
	}

	return map[string]cli.Command{
		createSubcommand.Name(): &createSubcommand,
		listSubcommand.Name():   &listSubcommand,
		buildSubcommand.Name():  &buildSubcommand,
		startSubcommand.Name():  &startSubcommand,
	}
}
