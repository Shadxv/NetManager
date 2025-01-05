package minecraft

import (
	"NetManager/internal/cli/model"
	"fmt"
)

func SetupWizard(printer model.Printer, serviceName string) {
	printer.Pause()
	var serviceType string
	fmt.Print("Service type (paper/velocity): ")
	_, err := fmt.Scan(&serviceType)
	if err != nil {
		return
	}
	printer.Resume()
}
