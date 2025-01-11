package servicesub

import (
	"NetManager/external/cli"
	"NetManager/external/service"
	"bytes"
	"github.com/olekukonko/tablewriter"
)

type ListSubcommand struct {
	Printer        cli.Printer
	ServiceManager service.Manager
}

func (cmd *ListSubcommand) Execute(args []string) {
	buf := new(bytes.Buffer)
	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"Name", "Type", "Status", "Image", "Version"})

	services := cmd.ServiceManager.Services()
	for _, s := range services {
		table.Append([]string{
			s.Name(),
			s.ServiceType(),
			s.Status(),
			s.ImageName(),
			s.CurrentVersion(),
		})
	}
	table.Render()
	cmd.Printer.Print("\n"+buf.String(), cmd.Printer.Service())
}

func (cmd *ListSubcommand) Name() string {
	return "list"
}

func (cmd *ListSubcommand) Description() string {
	return "Lists all services"
}

func (cmd *ListSubcommand) Subcommands() map[string]cli.Command {
	return map[string]cli.Command{}
}
