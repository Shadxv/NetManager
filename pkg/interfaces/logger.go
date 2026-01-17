package interfaces

type Logger interface {
	Init()
	IsPresent() bool
	Finalize() error
	ReadFile()
	ReadNext() (string, data []string)
	SaveLog(module string, data []string)
	Close()
}
