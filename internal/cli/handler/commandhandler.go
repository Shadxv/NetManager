package handler

import (
	"NetManager/external/cli"
	"strings"
)

func HandleHelp(printer cli.Printer, command cli.Command) {
	help := "Help for '" + command.Name() + "' - " + command.Description()
	if len(command.Subcommands()) > 0 {
		help += "\n\nSubcommands:"
		for subcmdName := range command.Subcommands() {
			help += "\n"
			subcmd := command.Subcommands()[subcmdName]
			help += "- " + subcmd.Name() + " - " + subcmd.Description()
		}
	}
	printer.Print("\n"+help, printer.Service())
}

func HandleSubcommand(printer cli.Printer, command cli.Command, args []string) {
	if len(args) <= 0 {
		HandleHelp(printer, command)
		return
	}
	subcommandName := strings.ToLower(args[0])
	subcommand := command.Subcommands()[subcommandName]
	if subcommand == nil {
		HandleHelp(printer, command)
		return
	}
	args = args[1:]
	subcommand.Execute(args)
}
