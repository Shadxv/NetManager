package manager

import (
	"NetManager/external/cli"
	"NetManager/internal/cli/commands"
	"strings"
)

type CommandManager struct {
	printer  cli.Printer
	commands map[string]cli.Command
}

func NewCommandManager(printer cli.Printer) *CommandManager {
	return &CommandManager{
		printer:  printer,
		commands: make(map[string]cli.Command),
	}
}

func (commandManager *CommandManager) ExecuteCommand(input string) {
	var split = strings.Split(input, " ")
	var commandName = strings.ToLower(split[0])
	var arguments = split[1:]

	command, exist := commandManager.commands[commandName]
	if !exist {
		commandManager.printer.PrintColored("Command does not exist!", commandManager.printer.Service(), cli.Red)
		return
	}

	command.Execute(arguments)
}

func (commandManager *CommandManager) RegisterCommand(command cli.Command) {
	commandManager.commands[command.Name()] = command
}

func (commandManager *CommandManager) registerCommands() {
	commandManager.RegisterCommand(&commands.StopCommand{
		Printer: commandManager.printer,
	})
}

func (commandManager *CommandManager) Init() {
	commandManager.registerCommands()
}
