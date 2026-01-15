package manager

import (
	"NetManager/internal/cli/commands"
	harborModel "NetManager/internal/docker/model"
	"NetManager/internal/kubernetes"
	"NetManager/internal/module"
	mongoModel "NetManager/internal/mongodb/model"
	redisModel "NetManager/internal/redis/model"
	"NetManager/internal/service/model"
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"NetManager/pkg/util"
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
	printer       interfaces.Printer       `gob:"-"`
	configManager interfaces.ConfigManager `gob:"-"`
	Services      map[string]*model.Service
	status        string `gob:"-"`
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		Services: make(map[string]*model.Service),
		status:   types.Stopped,
	}
}

func (manager *ServiceManager) Init(moduleManager *module.Manager) {
	manager.status = types.Starting
	printer, err := module.GetTypedModule[interfaces.Printer](moduleManager, types.Console)
	if err != nil {
		manager.status = types.Disabled
		return
	}
	manager.printer = printer

	configManager, err := module.GetTypedModule[interfaces.ConfigManager](moduleManager, types.Config)
	if err != nil {
		manager.printer.PrintColored("Config module not found!", manager.printer.Service(), types.Red)
		manager.status = types.Disabled
		return
	}
	manager.configManager = configManager

	k8sClient, err := module.GetTypedModule[*kubernetes.Client](moduleManager, types.Kubernetes)
	var clusterManager interfaces.ClusterManager
	if err != nil {
		manager.printer.PrintColored("Kubernetes module not found or not enabled!", manager.printer.Service(), types.Red)
	} else {
		clusterManager = k8sClient.ClusterManager()
	}

	if !manager.Exists("harbor") && clusterManager != nil {
		harborModel.CreateHarborService(manager.printer, manager.configManager.GetHarborConfig(), manager, clusterManager)
	}

	if !manager.Exists("redis") && clusterManager != nil {
		redisModel.CreateNewRedisService(manager.printer, manager.configManager.GetRedisConfig(), manager, clusterManager)
	}

	if !manager.Exists("mongodb") {
		mongoModel.CreateNewMongoService(manager.configManager.GetMongoConfig(), manager)
	}

	manager.printer.CommandManager().RegisterCommand(&commands.ServiceCommand{
		Printer:       manager.printer,
		ModuleManager: moduleManager,
	})
	manager.status = types.Enabled
}

func (manager *ServiceManager) Disable(shutdown bool) {
	if !shutdown {
		if manager.printer != nil {
			manager.printer.PrintColored("This module cannot be disabled!", manager.printer.Service(), types.Red)
		}
		return
	}
	manager.status = types.Stopping

	err := manager.SaveData()
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.printer.Service(), types.Red)
	}
	// stop watchers
	// clear Services
	manager.status = types.Disabled
}

func (manager *ServiceManager) Reload() {
	if manager.printer != nil {
		manager.printer.PrintColored("This module cannot be reloaded!", manager.printer.Service(), types.Red)
	}
}

func (manager *ServiceManager) SaveData() error {
	return util.SaveData(manager.Type(), manager)
}

func (manager *ServiceManager) LoadData() {
	if m, err := util.LoadData[ServiceManager](types.Services); err == nil {
		manager.Services = m.Services
	}
}

func (manager *ServiceManager) SetStatus(newStatus string) {
	manager.status = newStatus
}

func (manager *ServiceManager) Status() string {
	return manager.status
}

func (manager *ServiceManager) Type() string {
	return types.Services
}

func (manager *ServiceManager) AddNewService(name string, serviceType string, serviceData interface{}) interfaces.ServiceModel {
	serviceModel := model.CreateNewService(name, serviceType, serviceData)
	manager.Services[name] = serviceModel
	go manager.createMinecraftServiceFileStructure(serviceModel)
	return serviceModel
}

func (manager *ServiceManager) AddService(name string, serviceType string, status string, image string, version string, serviceData interface{}) interfaces.ServiceModel {
	serviceModel := model.NewService(
		name,
		serviceType,
		status,
		image,
		version,
		make([]string, 0),
		serviceData,
	)
	manager.Services[name] = serviceModel
	go manager.createMinecraftServiceFileStructure(serviceModel)
	return serviceModel
}

func (manager *ServiceManager) DeleteService(name string) interfaces.ServiceModel {
	serviceModel := manager.Services[name]
	if serviceModel == nil {
		return nil
	}
	delete(manager.Services, name)
	return serviceModel
}

func (manager *ServiceManager) Exists(name string) bool {
	return manager.Services[name] != nil
}

func (manager *ServiceManager) GetService(name string) interfaces.ServiceModel {
	return manager.Services[name]
}

func (manager *ServiceManager) GetServices() []interfaces.ServiceModel {
	services := make([]interfaces.ServiceModel, 0, len(manager.Services))
	for _, serviceModel := range manager.Services {
		services = append(services, serviceModel)
	}
	return services
}

// Minecraft section
func (manager *ServiceManager) createMinecraftServiceFileStructure(service interfaces.ServiceModel) {
	switch service.ServiceType() {
	case types.Paper, types.Velocity:
		manager.createFileStructure(service)
	default:
		return
	}
}

func (manager *ServiceManager) createFileStructure(service interfaces.ServiceModel) {
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

func (manager *ServiceManager) downloadServerJar(service interfaces.ServiceModel, jarPath string) bool {
	out, err := os.Create(jarPath)
	if err != nil {
		manager.handleStructureError(err, "Error occured during creating file structure for "+service.Name()+" service. Removing it from service registry.", service)
		return false
	}
	defer out.Close()

	url := manager.buidlUrlForServiceJar(service)
	if url == "" {
		manager.printer.PrintColored("Error occured during building url for server jar for "+service.Name()+" service. Removing it from service registry.", manager.printer.Service(), types.Red)
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
		manager.printer.PrintColored("Error occured during downloading server jar for "+service.Name()+" service. Removing it from service registry. Status code: "+resp.Status, manager.printer.Service(), types.Red)
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

func (manager *ServiceManager) handleStructureError(err error, msg string, service interfaces.ServiceModel) {
	manager.printer.PrintColored(msg, manager.printer.Service(), types.Red)
	manager.printer.PrintColored(err.Error(), manager.printer.Service(), types.Red)
	manager.DeleteService(service.Name())
}

func (manager *ServiceManager) buidlUrlForServiceJar(service interfaces.ServiceModel) string {
	switch service.ServiceType() {
	case types.Paper:
		data := interfaces.GetPaperData(service)
		if data == nil {
			manager.printer.PrintColored("Data is nil.", manager.printer.Service(), types.Red)
			return ""
		}
		version := data.Version()
		build := strconv.Itoa(data.BuildNumber())
		return "https://api.papermc.io/v2/projects/paper/versions/" + version + "/builds/" + build + "/downloads/paper-" + version + "-" + build + ".jar"
	case types.Velocity:
		data := interfaces.GetVelocityData(service)
		if data == nil {
			manager.printer.PrintColored("Data is nil.", manager.printer.Service(), types.Red)
			return ""
		}
		version := data.Version()
		build := strconv.Itoa(data.BuildNumber())
		return "https://api.papermc.io/v2/projects/velocity/versions/" + version + "/builds/" + build + "/downloads/velocity-" + version + "-" + build + ".jar"
	default:
		return ""
	}
}

func (manager *ServiceManager) createDockerfile(service interfaces.ServiceModel, path string) bool {
	dockerfileContent := fmt.Sprintf(`
FROM eclipse-temurin:21-jre
WORKDIR /dreammc
COPY %s-default/ /dreammc/
COPY %s/ /dreammc/
RUN echo "eula=true" > eula.txt
CMD ["java", "-jar", "-Xmx4G", "server.jar", "--nogui"]
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
