package model

type Command interface {
	Execute(args []string)
	Name() string
	Description() string
}
