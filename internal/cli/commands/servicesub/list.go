package servicesub

import (
	"NetManager/internal/cli/model"
	types "NetManager/internal/minecraft/type"
	"bytes"
	"github.com/olekukonko/tablewriter"
)

type ListSubcommand struct {
	Printer        model.Printer
	ServiceManager types.ServiceManagerBase
}

func (cmd *ListSubcommand) Execute(args []string) {
	buf := new(bytes.Buffer)
	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"Name", "Type"})

	services := cmd.ServiceManager.GetAllServices()
	for _, service := range services {
		table.Append([]string{
			service.GetName(),
			service.GetType(),
		})
	}
	table.Render()
	cmd.Printer.Print("\n" + buf.String(), cmd.Printer.Service())
}

func (cmd *ListSubcommand) Name() string {
	return "list"
}

func (cmd *ListSubcommand) Description() string {
	return "Lists all services"
}

func (cmd *ListSubcommand) Subcommands() map[string]model.Command {
	return map[string]model.Command{}
}
