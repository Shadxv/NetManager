package model

type Service interface {
	GetType() string
	GetName() string
	GetVersion() string
	GetBuild() int
}
