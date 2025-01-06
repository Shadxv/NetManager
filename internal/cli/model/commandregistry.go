package model

type CommandRegistry interface {
	RegisterCommand(command Command)
}
