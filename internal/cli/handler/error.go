package handler

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
)

func HandleError(printer interfaces.Printer, message string, err error, service interfaces.Service, shouldShutdown bool) bool {
	if err != nil {
		printer.PrintColored(message, service, types.Red)
		printer.PrintColored(err.Error(), service, types.Red)
		if shouldShutdown {
			printer.CloseGracefully("App is shutting down...")
		}
		return true
	}
	return false
}
