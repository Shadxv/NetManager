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
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

type ServiceManager struct {
	printer        interfaces.Printer        `gob:"-"`
	configManager  interfaces.ConfigManager  `gob:"-"`
	clusterManager interfaces.ClusterManager `gob:"-"`
	broadcaster    interfaces.Broadcaster    `gob:"-"`
	Services       map[string]*model.Service
	status         string `gob:"-"`
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		Services: make(map[string]*model.Service),
		status:   types.Stopped,
	}
}

func (manager *ServiceManager) SetBroadcaster(broadcaster interfaces.Broadcaster) {
	manager.broadcaster = broadcaster
	for _, service := range manager.Services {
		service.SetBroadcaster(broadcaster)
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
	if err != nil {
		manager.printer.PrintColored("Kubernetes module not found or not enabled!", manager.printer.Service(), types.Red)
	} else {
		manager.clusterManager = k8sClient.ClusterManager()
	}

	if !manager.Exists("harbor") && manager.clusterManager != nil {
		harborModel.CreateHarborService(manager.printer, manager.configManager.GetHarborConfig(), manager, manager.clusterManager)
	}

	if !manager.Exists("redis") && manager.clusterManager != nil {
		redisModel.CreateNewRedisService(manager.printer, manager.configManager.GetRedisConfig(), manager, manager.clusterManager)
	} else {
		manager.GetService("redis").Deploy(manager.printer, manager.clusterManager)
	}

	if !manager.Exists("mongodb") {
		mongoModel.CreateNewMongoService(manager.configManager.GetMongoConfig(), manager)
	}

	manager.printer.CommandManager().RegisterCommand(&commands.ServiceCommand{
		Printer:       manager.printer,
		ModuleManager: moduleManager,
	})

	if manager.clusterManager != nil {
		manager.syncPods()
		go manager.watchPods()
	}

	manager.status = types.Enabled
}

func (manager *ServiceManager) syncPods() {
	pods := manager.clusterManager.GetPods("", manager.clusterManager.GetDefaultNamespace())
	for _, pod := range pods {
		manager.handlePodUpdate(&pod, "ADDED")
	}
}

func (manager *ServiceManager) watchPods() {
	watcher, err := manager.clusterManager.WatchPods("", manager.clusterManager.GetDefaultNamespace())
	if err != nil {
		manager.printer.PrintColored("Error watching pods: "+err.Error(), manager.printer.Service(), types.Red)
		return
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		pod, ok := event.Object.(*corev1.Pod)
		if !ok {
			continue
		}
		manager.handlePodUpdate(pod, string(event.Type))
	}
}

func (manager *ServiceManager) handlePodUpdate(pod *corev1.Pod, eventType string) {
	serviceName := pod.Labels["app"]
	if serviceName == "" {
		return
	}

	service := manager.GetService(serviceName)
	if service == nil {
		return
	}
	if service.PodInstances() == nil {
		return
	}

	switch eventType {
	case "ADDED", "MODIFIED":
		var instance interfaces.PodInstance
		if existingInstance, ok := service.PodInstances()[pod.Name]; ok {
			instance = existingInstance
			oldStatus := instance.Status()
			instance.UpdatePod(pod)
			instance.SetStatus(string(pod.Status.Phase))

			if oldStatus != "Running" && instance.Status() == "Running" {
				go manager.streamLogs(instance, serviceName)
			}
		} else {
			instance = service.AddPodInstance(pod.Name, pod, string(pod.Status.Phase))
			go manager.streamLogs(instance, serviceName)
		}

		instance.SetInternalIP(pod.Status.PodIP)
		if pod.Status.HostIP != "" {
			instance.SetExternalIP(pod.Status.HostIP)
		}

		if len(pod.Spec.Containers) > 0 && len(pod.Spec.Containers[0].Ports) > 0 {
			instance.SetPort(pod.Spec.Containers[0].Ports[0].ContainerPort)
		}

		if manager.broadcaster != nil {
			manager.broadcaster.BroadcastEvent(serviceName, "pod_update", map[string]interface{}{
				"pod":    pod.Name,
				"status": string(pod.Status.Phase),
				"event":  eventType,
			})
		}

	case "DELETED":
		service.RemovePodInstance(pod.Name)
		if manager.broadcaster != nil {
			manager.broadcaster.BroadcastEvent(serviceName, "pod_deleted", map[string]interface{}{
				"pod": pod.Name,
			})
		}
	}
}

func (manager *ServiceManager) streamLogs(instance interfaces.PodInstance, serviceName string) {
	logs, err := manager.clusterManager.GetPodLogs(instance.Name(), manager.clusterManager.GetDefaultNamespace())
	if err != nil {
		return
	}
	defer logs.Close()

	scanner := bufio.NewScanner(logs)
	for scanner.Scan() {
		line := scanner.Text()
		instance.AddLog(line)
		if manager.broadcaster != nil {
			manager.broadcaster.BroadcastLog(serviceName, instance.Name(), line)
		}
	}
}

func (manager *ServiceManager) Disable(shutdown bool) {
	if !shutdown {
		if manager.printer != nil {
			manager.printer.PrintColored("This module cannot be disabled!", manager.printer.Service(), types.Red)
		}
		return
	}
	manager.status = types.Stopping
	manager.SaveData()
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
	serviceModel := model.CreateNewService(name, serviceType, manager.clusterManager.GetDefaultNamespace(), serviceData)
	serviceModel.SetStatus(types.Creating)
	manager.Services[name] = serviceModel
	go manager.createMinecraftServiceFileStructure(serviceModel)
	if manager.broadcaster != nil {
		serviceModel.SetBroadcaster(manager.broadcaster)
	}
	return serviceModel
}

func (manager *ServiceManager) AddService(name string, serviceType string, status string, image string, namespace string, version string, serviceData interface{}) interfaces.ServiceModel {
	serviceModel := model.NewService(
		name,
		serviceType,
		status,
		image,
		namespace,
		version,
		make([]string, 0),
		serviceData,
	)
	manager.Services[name] = serviceModel
	go manager.createMinecraftServiceFileStructure(serviceModel)
	if manager.broadcaster != nil {
		serviceModel.SetBroadcaster(manager.broadcaster)
	}
	return serviceModel
}

func (manager *ServiceManager) DeleteService(name string) interfaces.ServiceModel {
	if name == "mongodb" || name == "redis" || name == "harbor" {
		manager.printer.PrintColored("Cannot delete protected service: "+name, manager.printer.Service(), types.Red)
		return nil
	}

	serviceModel := manager.Services[name]
	if serviceModel == nil {
		return nil
	}

	if manager.clusterManager != nil {
		serviceModel.Stop(manager.printer, manager.clusterManager)
	}

	folderPath := filepath.Join("services/templates", name)
	err := os.RemoveAll(folderPath)
	if err != nil {
		manager.printer.PrintColored("Error deleting service files: "+err.Error(), manager.printer.Service(), types.Red)
	}

	delete(manager.Services, name)
	manager.printer.Print("Service "+name+" deleted successfully.", manager.printer.Service())

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
	service.SetStatus(types.Stopped)
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

func (manager *ServiceManager) buildRunCommand(service interfaces.ServiceModel) string {
	var cmd string
	switch service.ServiceType() {
	case types.Paper:
		cmd = `CMD ["java", "-jar", "-Xmx4G", "server.jar", "--nogui"]`
		break
	case types.Velocity:
		cmd = `CMD ["java", "-jar", "-Xmx4G", "server.jar"]`
		break
	default:
		return ""
	}
	return cmd
}

func (manager *ServiceManager) createDockerfile(service interfaces.ServiceModel, path string) bool {
	dockerfileContent := fmt.Sprintf(`
FROM eclipse-temurin:21-jre
WORKDIR /dreammc
COPY %s-default/ /dreammc/
COPY %s/ /dreammc/
RUN echo "eula=true" > eula.txt
%s
`, strings.ToLower(service.ServiceType()), service.Name(), manager.buildRunCommand(service))

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
