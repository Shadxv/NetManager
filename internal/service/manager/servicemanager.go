package manager

import (
	"NetManager/external/cli"
	"NetManager/external/docker"
	"NetManager/external/kubernetes"
	"NetManager/external/minecraft"
	"NetManager/external/service"
	"NetManager/external/types"
	"NetManager/internal/cli/commands"
	"NetManager/internal/service/model"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type ServiceManager struct {
	printer  cli.Printer
	services map[string]service.Model
}

func NewServiceManager(printer cli.Printer, services map[string]service.Model) *ServiceManager {
	return &ServiceManager{
		printer:  printer,
		services: services,
	}
}

func CreateNewServiceManager(printer cli.Printer) *ServiceManager {
	return &ServiceManager{
		printer:  printer,
		services: make(map[string]service.Model),
	}
}

func (manager *ServiceManager) Init(commandManager cli.CommandManager, imageManager docker.ImageManager, kubernetesClient kubernetes.Client) *ServiceManager {
	commandManager.RegisterCommand(&commands.ServiceCommand{
		Printer:          manager.printer,
		ServiceManager:   manager,
		ImageManager:     imageManager,
		KubernetesClient: kubernetesClient,
	})
	return manager
}

func (manager *ServiceManager) AddNewService(name string, serviceType string, serviceData *interface{}) service.Model {
	serviceModel := model.CreateNewService(name, serviceType, serviceData)
	manager.services[name] = serviceModel
	go manager.createMinecraftServiceFileStructure(serviceModel)
	return serviceModel
}

func (manager *ServiceManager) AddService(name string, serviceType string, status string, image string, version string, serviceData interface{}) service.Model {
	serviceModel := model.NewService(
		name,
		serviceType,
		status,
		image,
		version,
		make([]string, 0),
		&serviceData,
	)
	manager.services[name] = serviceModel
	go manager.createMinecraftServiceFileStructure(serviceModel)
	return serviceModel
}

func (manager *ServiceManager) DeleteService(name string) service.Model {
	serviceModel := manager.services[name]
	if serviceModel == nil {
		return nil
	}
	delete(manager.services, name)
	return serviceModel
}

func (manager *ServiceManager) Exists(name string) bool {
	return manager.services[name] != nil
}

func (manager *ServiceManager) GetService(name string) service.Model {
	return manager.services[name]
}

func (manager *ServiceManager) Services() []service.Model {
	services := make([]service.Model, 0, len(manager.services))
	for _, serviceModel := range manager.services {
		services = append(services, serviceModel)
	}
	return services
}

// Minecraft section
func (manager *ServiceManager) createMinecraftServiceFileStructure(service service.Model) {
	switch service.ServiceType() {
	case types.Paper, types.Velocity:
		manager.createFileStructure(service)
	default:
		return
	}
}

func (manager *ServiceManager) createFileStructure(service service.Model) {
	manager.printer.Print("Creating file structure for "+service.Name()+"...", manager.printer.Service())
	folderPath := filepath.Join("services/templates", service.Name())
	err := os.Mkdir(folderPath, 0755)
	if err != nil {
		manager.handleStructureError(err, "Error occured during creating file structure for "+service.Name()+" service. Removing it from service registry.", service)
		return
	}

	manager.printer.Print("Starting downloading server jar for "+service.Name()+"...", manager.printer.Service())
	jarPath := filepath.Join(folderPath, "server.jar")
	if !manager.downloadServerJar(service, jarPath) || !manager.createDockerfile(service, path.Join(folderPath, "Dockerfile")) {
		return
	}
	manager.printer.Print("Successfuly downloaded server jar for "+service.Name()+".", manager.printer.Service())
	manager.printer.Print("Service "+service.Name()+" creation has finished.", manager.printer.Service())
}

func (manager *ServiceManager) downloadServerJar(service service.Model, jarPath string) bool {
	out, err := os.Create(jarPath)
	if err != nil {
		manager.handleStructureError(err, "Error occured during creating file structure for "+service.Name()+" service. Removing it from service registry.", service)
		return false
	}
	defer out.Close()

	url := manager.buidlUrlForServiceJar(service)
	if url == "" {
		manager.printer.PrintColored("Error occured during building url for server jar for "+service.Name()+" service. Removing it from service registry.", manager.printer.Service(), cli.Red)
		manager.DeleteService(service.Name())
		return false
	}

	resp, err := http.Get(url)
	if err != nil {
		manager.handleStructureError(err, "Error occured during downloading server jar for "+service.Name()+" service. Removing it from service registry.", service)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		manager.printer.PrintColored("Error occured during downloading server jar for "+service.Name()+" service. Removing it from service registry. Status code: "+resp.Status, manager.printer.Service(), cli.Red)
		manager.DeleteService(service.Name())
		return false
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		manager.handleStructureError(err, "Error occured during saving server jar for "+service.Name()+" service. Removing it from service registry.", service)
		return false
	}
	return true
}

func (manager *ServiceManager) handleStructureError(err error, msg string, service service.Model) {
	manager.printer.PrintColored(msg, manager.printer.Service(), cli.Red)
	manager.printer.PrintColored(err.Error(), manager.printer.Service(), cli.Red)
	manager.DeleteService(service.Name())
}

func (manager *ServiceManager) buidlUrlForServiceJar(service service.Model) string {
	switch service.ServiceType() {
	case types.Paper:
		data := *minecraft.GetPaperData(service)
		version := data.Version()
		build := strconv.Itoa(data.Build())
		return "https://api.papermc.io/v2/projects/paper/versions/" + version + "/builds/" + build + "/downloads/paper-" + version + "-" + build + ".jar"
	case types.Velocity:
		data := *minecraft.GetVelocityData(service)
		version := data.Version()
		build := strconv.Itoa(data.Build())
		return "https://api.papermc.io/v2/projects/velocity/versions/" + version + "/builds/" + build + "/downloads/velocity-" + version + "-" + build + ".jar"
	default:
		return ""
	}
}

func (manager *ServiceManager) createDockerfile(service service.Model, path string) bool {
	dockerfileContent := fmt.Sprintf(`
FROM eclipse-temurin:21-jre
WORKDIR /dreammc
COPY %s-default/ /dreammc/
COPY %s/ /dreammc/
RUN echo "eula=true" > eula.txt
CMD ["java -jar server.jar nogui"]
`, strings.ToLower(service.ServiceType()), service.Name())

	out, err := os.Create(path)
	if err != nil {
		manager.handleStructureError(err, "Error occured during creating Dockerfile for "+service.Name()+" service. Removing it from service registry.", service)
		return false
	}
	defer out.Close()

	_, err = out.WriteString(dockerfileContent)
	if err != nil {
		manager.handleStructureError(err, "Error occured during creating Dockerfile for "+service.Name()+" service. Removing it from service registry.", service)
		os.Remove(path)
		return false
	}

	return true
}
