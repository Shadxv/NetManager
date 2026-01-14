package interfaces

type Printer interface {
	PrintColored(message string, service Service, color int)
	Print(message string, service Service)
	SetClosingStatus()
	Close()
	CloseGracefully(message string)
	Pause()
	Resume()
	Service() Service
	CommandManager() CommandManager
}
