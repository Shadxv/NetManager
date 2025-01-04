package minecraft

import (
	"NetManager/internal/cli/model"
	"fmt"
)

func SetupWizard(printer model.Printer, serviceName string) {
	printer.Pause()
	printer.Print("Paused...", printer.Service())
	var serviceType string
	fmt.Print("Service type (paper/velocity): ")
	fmt.Scan(&serviceType)
	printer.Print("Resumed...", printer.Service())
	printer.Resume()
}
