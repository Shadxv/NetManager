package servicesub

import (
	"NetManager/external/cli"
	"NetManager/external/kubernetes"
	"NetManager/external/service"

	"strings"
)

type StartSubcommand struct {
	Printer          cli.Printer
	ServiceManager   service.Manager
	KubernetesClient kubernetes.Client
}

func (cmd *StartSubcommand) Execute(args []string) {
	if len(args) <= 0 {
		cmd.Printer.PrintColored("Service name not specified. 'service create <name>'", cmd.Printer.Service(), cli.Yellow)
		return
	} else if len(args) > 1 {
		cmd.Printer.PrintColored("Service name cannot contain spaces.", cmd.Printer.Service(), cli.Yellow)
		return
	}

	name := strings.ToLower(args[0])
	if !cmd.ServiceManager.Exists(name) {
		cmd.Printer.PrintColored("Service with this name does not exist.", cmd.Printer.Service(), cli.Yellow)
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

func (cmd *StartSubcommand) Subcommands() map[string]cli.Command {
	return map[string]cli.Command{}
}
