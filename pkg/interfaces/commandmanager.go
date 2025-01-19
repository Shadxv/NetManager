package interfaces

type CommandManager interface {
	Init()
	ExecuteCommand(input string)
	RegisterCommand(command Command)
}
