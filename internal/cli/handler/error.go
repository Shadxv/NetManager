package handler

import (
	"NetManager/external/cli"
)

func HandleError(printer cli.Printer, message string, err error, service cli.Service, shouldShutdown bool) bool {
	if err != nil {
		printer.PrintColored(message, service, cli.Red)
		printer.PrintColored(err.Error(), service, cli.Red)
		if shouldShutdown {
			printer.CloseGracefully("App is shutting down...")
		}
		return true
	}
	return false
}
