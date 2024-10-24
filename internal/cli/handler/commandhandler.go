package handler

import (
	"NetManager/internal/cli/model"
	"strings"
)

func HandleHelp(printer model.Printer, command model.Command) {
	printer.Print("Help", printer.Service())
}

func HandleSubcommand(printer model.Printer, command model.Command, args []string) {
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
