package servicesub

import (
	"NetManager/external/cli"
	"NetManager/external/service"
	"bytes"
	"github.com/olekukonko/tablewriter"
)

type ListPodsSubcommand struct {
	Printer        cli.Printer
	ServiceManager service.Manager
}

func (cmd *ListPodsSubcommand) Execute(args []string) {

	if len(args) != 1 {
		cmd.Printer.PrintColored("Service name not specified. 'service listpods <name>'", cmd.Printer.Service(), cli.Yellow)
		return
	}

	buf := new(bytes.Buffer)
	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"Name", "Status"})

	serviceModel := cmd.ServiceManager.GetService(args[0])
	if serviceModel == nil {
		cmd.Printer.PrintColored("Service not found.", cmd.Printer.Service(), cli.Yellow)
		return
	}

	for _, p := range serviceModel.PodInstances() {
		table.Append([]string{
			p.Name(),
			p.Status(),
		})
	}
	table.Render()
	cmd.Printer.Print("\n"+buf.String(), cmd.Printer.Service())
}

func (cmd *ListPodsSubcommand) Name() string {
	return "listpods"
}

func (cmd *ListPodsSubcommand) Description() string {
	return "Lists all pods that currently are running"
}

func (cmd *ListPodsSubcommand) Subcommands() map[string]cli.Command {
	return map[string]cli.Command{}
}
