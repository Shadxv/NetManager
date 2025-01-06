package manager

import (
	"NetManager/internal/cli/commands"
	"NetManager/internal/cli/model"
	serviceModel "NetManager/internal/minecraft/model"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type ServiceManager struct {
	printer         model.Printer
	services        map[string]serviceModel.Service
	commandRegistry model.CommandRegistry
}

func NewServiceManager(printer model.Printer, commandRegistry model.CommandRegistry) *ServiceManager {
	return &ServiceManager{
		printer:         printer,
		services:        make(map[string]serviceModel.Service),
		commandRegistry: commandRegistry,
	}
}

func (manager *ServiceManager) Init() {
	manager.commandRegistry.RegisterCommand(&commands.ServiceCommand{
		Printer:        manager.printer,
		ServiceManager: manager,
	})
}

func (manager *ServiceManager) GetAllServices() map[string]serviceModel.Service {
	return manager.services
}

func (manager *ServiceManager) RegisterNewService(service serviceModel.Service) {
	manager.services[service.GetName()] = service
	manager.printer.Print("Registered new service: "+service.GetName(), manager.printer.Service())
	go manager.createFileStructure(service)
}

func (manager *ServiceManager) createFileStructure(service serviceModel.Service) {
	manager.printer.Print("Creating file structure for "+service.GetName()+"...", manager.printer.Service())
	folderPath := filepath.Join("services/templates", service.GetName())
	err := os.Mkdir(folderPath, 0755)
	if err != nil {
		manager.handleStructureError(err, "Error occured during creating file structure for "+service.GetName()+" service. Removing it from service registry.", service)
		return
	}

	manager.printer.Print("Starting downloading server jar for "+service.GetName()+"...", manager.printer.Service())
	jarPath := filepath.Join(folderPath, "server.jar")
	if !manager.downloadServerJar(service, jarPath) {
		return
	}
	manager.printer.Print("Successfuly downloaded server jar for "+service.GetName()+".", manager.printer.Service())
	manager.printer.Print("Service "+service.GetName()+" creation has finished.", manager.printer.Service())
}

func (manager *ServiceManager) downloadServerJar(service serviceModel.Service, jarPath string) bool {
	out, err := os.Create(jarPath)
	if err != nil {
		manager.handleStructureError(err, "Error occured during creating file structure for "+service.GetName()+" service. Removing it from service registry.", service)
		return false
	}
	defer out.Close()

	url := manager.buidlUrlForServiceJar(service)
	if url == "" {
		manager.printer.PrintColored("Error occured during building url for server jar for "+service.GetName()+" service. Removing it from service registry.", manager.printer.Service(), model.Red)
		manager.clearServiceData(service)
		return false
	}

	resp, err := http.Get(url)
	if err != nil {
		manager.handleStructureError(err, "Error occured during downloading server jar for "+service.GetName()+" service. Removing it from service registry.", service)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		manager.printer.PrintColored("Error occured during downloading server jar for "+service.GetName()+" service. Removing it from service registry. Status code: "+resp.Status, manager.printer.Service(), model.Red)
		manager.clearServiceData(service)
		return false
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		manager.handleStructureError(err, "Error occured during saving server jar for "+service.GetName()+" service. Removing it from service registry.", service)
		return false
	}
	return true
}

func (manager *ServiceManager) clearServiceData(service serviceModel.Service) {
	delete(manager.services, service.GetName())
}

func (manager *ServiceManager) handleStructureError(err error, msg string, service serviceModel.Service) {
	manager.printer.PrintColored(msg, manager.printer.Service(), model.Red)
	manager.printer.PrintColored(err.Error(), manager.printer.Service(), model.Red)
	manager.clearServiceData(service)
}

func (manager *ServiceManager) buidlUrlForServiceJar(service serviceModel.Service) string {
	build := strconv.Itoa(service.GetBuild())
	manager.printer.Print(service.GetType()+" "+service.GetVersion()+" "+build, manager.printer.Service())
	switch service.GetType() {
	case serviceModel.PaperType:
		return "https://api.papermc.io/v2/projects/paper/versions/" + service.GetVersion() + "/builds/" + build + "/downloads/paper-" + service.GetVersion() + "-" + build + ".jar"
	case serviceModel.VelocityType:
		return "https://api.papermc.io/v2/projects/velocity/versions/" + service.GetVersion() + "/builds/" + build + "/downloads/velocity-" + service.GetVersion() + "-" + build + ".jar"
	default:
		return ""
	}
}
