package servicesub

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"

	"strings"
)

type StopSubcommand struct {
	Printer          interfaces.Printer
	ServiceManager   interfaces.ServiceManager
	KubernetesClient interfaces.KubernetesClient
}

func (cmd *StopSubcommand) Execute(args []string) {
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

	cmd.ServiceManager.GetService(name).Stop(cmd.Printer, cmd.KubernetesClient.ClusterManager())
}

func (cmd *StopSubcommand) Name() string {
	return "stop"
}

func (cmd *StopSubcommand) Description() string {
	return "Stops service"
}

func (cmd *StopSubcommand) Subcommands() map[string]interfaces.Command {
	return map[string]interfaces.Command{}
}
