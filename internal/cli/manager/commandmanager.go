package manager

import (
	"NetManager/internal/cli/commands"
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"strings"
)

type CommandManager struct {
	printer  interfaces.Printer
	commands map[string]interfaces.Command
}

func NewCommandManager(printer interfaces.Printer) *CommandManager {
	return &CommandManager{
		printer:  printer,
		commands: make(map[string]interfaces.Command),
	}
}

func (commandManager *CommandManager) ExecuteCommand(input string) {
	var split = strings.Split(input, " ")
	var commandName = strings.ToLower(split[0])
	var arguments = split[1:]

	command, exist := commandManager.commands[commandName]
	if !exist {
		commandManager.printer.PrintColored("Command does not exist!", commandManager.printer.Service(), types.Red)
		return
	}

	command.Execute(arguments)
}

func (commandManager *CommandManager) RegisterCommand(command interfaces.Command) {
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
