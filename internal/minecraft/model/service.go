package model

var PaperType = "PAPER"
var VelocityType = "VELOCITY"

type Service interface {
	GetType() string
	GetName() string
	GetVersion() string
	GetBuild() int
}
