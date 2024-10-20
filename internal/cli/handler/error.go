package handler

import (
	"NetManager/internal/cli/model"
	"NetManager/internal/service"
)

func HandleError(printer model.Printer, message string, err error, service service.Service, shouldShutdown bool) bool {
	if err != nil {
		printer.PrintColored(message, service, model.Red)
		printer.PrintColored(err.Error(), service, model.Red)
		if shouldShutdown {
			printer.CloseGracefully("App is shutting down...")
		}
		return true
	}
	return false
}
