package cli

const (
	Red    = 160
	Green  = 72
	Yellow = 221
)

type Printer interface {
	PrintColored(message string, service Service, color int)
	Print(message string, service Service)
	SetClosingStatus()
	Close()
	CloseGracefully(message string)
	Pause()
	Resume()
	Service() Service
}
