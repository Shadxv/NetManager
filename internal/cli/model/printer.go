package model

import "NetManager/internal/service"

const (
	Red    = 160
	Green  = 72
	Yellow = 221
)

type Printer interface {
	PrintColored(message string, service service.Service, color int)
	Print(message string, service service.Service)
	SetClosingStatus()
	Close()
	CloseGracefully(message string)
	Pause()
	Resume()
	Service() service.Service
}
