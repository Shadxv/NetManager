package manager

import (
	"NetManager/internal/cli/commands"
	"NetManager/internal/cli/model"

	"strings"
)

type CommandManager struct {
	printer  model.Printer
	commands map[string]model.Command
}

func NewCommandManager(printer model.Printer) *CommandManager {
	return &CommandManager{
		printer:  printer,
		commands: make(map[string]model.Command),
	}
}

func (commandManager *CommandManager) ExecuteCommand(input string) {
	var split = strings.Split(input, " ")
	var commandName = strings.ToLower(split[0])
	var arguments = split[1:]

	command, exist := commandManager.commands[commandName]
	if !exist {
		commandManager.printer.PrintColored("Command does not exist!", commandManager.printer.Service(), model.Red)
		return
	}

	command.Execute(arguments)
}

func (commandManager *CommandManager) RegisterCommand(command model.Command) {
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
