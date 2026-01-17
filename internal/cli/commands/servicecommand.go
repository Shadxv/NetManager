package commands

import (
	"NetManager/internal/cli/commands/servicesub"
	"NetManager/internal/cli/handler"
	"NetManager/internal/kubernetes"
	"NetManager/internal/module"
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
)

type ServiceCommand struct {
	Printer       interfaces.Printer
	ModuleManager *module.Manager
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
	serviceManager, err := module.GetTypedModule[interfaces.ServiceManager](cmd.ModuleManager, types.Services)
	if err != nil {
		cmd.Printer.PrintColored("ServiceManager module not found or not enabled!", cmd.Printer.Service(), types.Red)
		return nil
	}

	configManager, err := module.GetTypedModule[interfaces.ConfigManager](cmd.ModuleManager, types.Config)
	if err != nil {
		cmd.Printer.PrintColored("ConfigManager module not found or not enabled!", cmd.Printer.Service(), types.Red)
		return nil
	}

	imageManager, err := module.GetTypedModule[interfaces.ImageManager](cmd.ModuleManager, types.Images)
	if err != nil {
		cmd.Printer.PrintColored("ImageManager module not found or not enabled!", cmd.Printer.Service(), types.Red)
		return nil
	}

	kubernetesClient, err := module.GetTypedModule[*kubernetes.Client](cmd.ModuleManager, types.Kubernetes)
	if err != nil {
		cmd.Printer.PrintColored("KubernetesClient module not found or not enabled!", cmd.Printer.Service(), types.Red)
		return nil
	}

	createSubcommand := servicesub.CreateSubcommand{
		Printer:        cmd.Printer,
		ConfigManager:  configManager,
		ServiceManager: serviceManager,
	}

	listSubcommand := servicesub.ListSubcommand{
		Printer:        cmd.Printer,
		ServiceManager: serviceManager,
	}

	listPodsSubcommand := servicesub.ListPodsSubcommand{
		Printer:        cmd.Printer,
		ServiceManager: serviceManager,
	}

	buildSubcommand := servicesub.BuildSubcommand{
		Printer:        cmd.Printer,
		ServiceManager: serviceManager,
		ImageManager:   imageManager,
	}

	startSubcommand := servicesub.StartSubcommand{
		Printer:          cmd.Printer,
		ServiceManager:   serviceManager,
		KubernetesClient: kubernetesClient,
	}

	stopSubcommand := servicesub.StopSubcommand{
		Printer:          cmd.Printer,
		ServiceManager:   serviceManager,
		KubernetesClient: kubernetesClient,
	}

	return map[string]interfaces.Command{
		createSubcommand.Name():   &createSubcommand,
		listSubcommand.Name():     &listSubcommand,
		buildSubcommand.Name():    &buildSubcommand,
		startSubcommand.Name():    &startSubcommand,
		stopSubcommand.Name():     &stopSubcommand,
		listPodsSubcommand.Name(): &listPodsSubcommand,
	}
}
