package commands

import (
	"NetManager/internal/cli/commands/servicesub"
	"NetManager/internal/cli/handler"
	"NetManager/pkg/interfaces"
)

type ServiceCommand struct {
	Printer          interfaces.Printer
	ServiceManager   interfaces.ServiceManager
	ImageManager     interfaces.ImageManager
	KubernetesClient interfaces.KubernetesClient
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

func (cmd *ServiceCommand) Subcommands() map[string]interfaces.Command {
	createSubcommand := servicesub.CreateSubcommand{
		Printer:        cmd.Printer,
		ServiceManager: cmd.ServiceManager,
	}

	listSubcommand := servicesub.ListSubcommand{
		Printer:        cmd.Printer,
		ServiceManager: cmd.ServiceManager,
	}

	listPodsSubcommand := servicesub.ListPodsSubcommand{
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

	return map[string]interfaces.Command{
		createSubcommand.Name():   &createSubcommand,
		listSubcommand.Name():     &listSubcommand,
		buildSubcommand.Name():    &buildSubcommand,
		startSubcommand.Name():    &startSubcommand,
		listPodsSubcommand.Name(): &listPodsSubcommand,
	}
}
