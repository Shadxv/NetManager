package servicesub

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"

	"strings"
)

type StartSubcommand struct {
	Printer          interfaces.Printer
	ServiceManager   interfaces.ServiceManager
	KubernetesClient interfaces.KubernetesClient
}

func (cmd *StartSubcommand) Execute(args []string) {
	if len(args) <= 0 {
		cmd.Printer.PrintColored("Service name not specified. 'service create <name>'", cmd.Printer.Service(), types.Yellow)
		return
	} else if len(args) > 1 {
		cmd.Printer.PrintColored("Service name cannot contain spaces.", cmd.Printer.Service(), types.Yellow)
		return
	}

	name := strings.ToLower(args[0])
	if !cmd.ServiceManager.Exists(name) {
		cmd.Printer.PrintColored("Service with this name does not exist.", cmd.Printer.Service(), types.Yellow)
		return
	}

	//cmd.KubernetesClient.DeployPaperService(name)
}

func (cmd *StartSubcommand) Name() string {
	return "start"
}

func (cmd *StartSubcommand) Description() string {
	return "Starts service"
}

func (cmd *StartSubcommand) Subcommands() map[string]interfaces.Command {
	return map[string]interfaces.Command{}
}
