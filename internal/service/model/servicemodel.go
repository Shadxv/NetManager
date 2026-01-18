package model

import (
	"NetManager/internal/kubernetes/model"
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"slices"

	corev1 "k8s.io/api/core/v1"
)

type Service struct {
	NameField              string
	ServiceTypeField       string
	StatusField            string
	ImageNameField         string
	NamespaceField         string
	CurrentVersionField    string
	AvailableVersionsField []string
	// Config from k8s
	appConfigField    interface{}                       `gob:"-"`
	podInstancesField map[string]interfaces.PodInstance `gob:"-"`
	netServiceField   *corev1.Service                   `gob:"-"`
	ServiceDataField  interface{}
	broadcaster       interfaces.Broadcaster `gob:"-"`
}

func NewService(name string, serviceType string, status string, imageName string, namespace string, currentVersion string, availableVersions []string, serviceData interface{}) *Service {
	return &Service{
		NameField:              name,
		ServiceTypeField:       serviceType,
		StatusField:            status,
		ImageNameField:         imageName,
		NamespaceField:         namespace,
		CurrentVersionField:    currentVersion,
		AvailableVersionsField: availableVersions,
		podInstancesField:      make(map[string]interfaces.PodInstance),
		ServiceDataField:       serviceData,
	}
}

func CreateNewService(name string, serviceType string, namespace string, serviceData interface{}) *Service {
	return &Service{
		NameField:         name,
		ServiceTypeField:  serviceType,
		StatusField:       types.Stopped,
		ImageNameField:    name,
		NamespaceField:    namespace,
		podInstancesField: make(map[string]interfaces.PodInstance),
		ServiceDataField:  serviceData,
	}
}

func (service *Service) InitNonSavedFields() {
	if service.podInstancesField == nil {
		service.podInstancesField = make(map[string]interfaces.PodInstance)
	}
}

func (service *Service) SetBroadcaster(broadcaster interfaces.Broadcaster) {
	service.broadcaster = broadcaster
}

func (service *Service) Name() string {
	return service.NameField
}

func (service *Service) ServiceType() string {
	return service.ServiceTypeField
}

func (service *Service) Status() string {
	return service.StatusField
}

func (service *Service) SetStatus(status string) {
	service.StatusField = status
	if service.broadcaster != nil {
		service.broadcaster.BroadcastEvent(service.Name(), "status", map[string]interface{}{
			"status": status,
		})
	}
}

func (service *Service) ImageName() string {
	return service.ImageNameField
}

func (service *Service) Namespace() string { return service.NamespaceField }

func (service *Service) CurrentVersion() string {
	return service.CurrentVersionField
}

func (service *Service) SetCurrentVersion(currentVersion string) bool {
	if slices.Contains(service.AvailableVersionsField, currentVersion) {
		service.CurrentVersionField = currentVersion
		if service.broadcaster != nil {
			service.broadcaster.BroadcastEvent(service.Name(), "current_version", map[string]interface{}{
				"version": currentVersion,
			})
		}
		return true
	}
	return false
}

func (service *Service) AvailableVersions() []string {
	versionsCopy := make([]string, len(service.AvailableVersionsField))
	copy(versionsCopy, service.AvailableVersionsField)
	return versionsCopy
}

func (service *Service) AddVersion(version string) {
	if slices.Contains(service.AvailableVersionsField, version) {
		return
	}
	service.AvailableVersionsField = append(service.AvailableVersionsField, version)
	if service.broadcaster != nil {
		service.broadcaster.BroadcastEvent(service.Name(), "version", map[string]interface{}{
			"newVersion": version,
		})
	}
}

func (service *Service) AppConfig() interface{} {
	return service.appConfigField
}

func (service *Service) SetAppConfig(config interface{}) {
	service.appConfigField = config
}

func (service *Service) RemoveAppConfig() {
	service.appConfigField = nil
}

func (service *Service) PodInstances() map[string]interfaces.PodInstance {
	return service.podInstancesField
}

func (service *Service) addPodInstance(instance interfaces.PodInstance) {
	service.podInstancesField[instance.Name()] = instance
}

func (service *Service) AddPodInstance(name string, pod *corev1.Pod, status string) interfaces.PodInstance {
	instance := model.NewPodInstance(name, pod, status)
	service.addPodInstance(instance)
	return instance
}

func (service *Service) CreatePodInstance(name string, pod *corev1.Pod) interfaces.PodInstance {
	instance := model.CreateNewPodInstance(name, pod)
	service.addPodInstance(instance)
	return instance
}

func (service *Service) RemovePodInstance(name string) {
	delete(service.podInstancesField, name)
}

func (service *Service) NetService() *corev1.Service {
	return service.netServiceField
}

func (service *Service) SetNetSevice(netService *corev1.Service) {
	service.netServiceField = netService
}

func (service *Service) RemoveNetService() {
	service.netServiceField = nil
}

func (service *Service) ServiceData() interface{} {
	return service.ServiceDataField
}

func (service *Service) Build(printer interfaces.Printer, imageManager interfaces.ImageManager, serviceManager interfaces.ServiceManager) {
	data := service.ServiceDataField
	data.(interfaces.Data).Build(service, printer, imageManager, serviceManager)
}

func (service *Service) Update(printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	data := service.ServiceDataField
	data.(interfaces.Data).Update(service, printer, clusterManager)
}

func (service *Service) Stop(printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	if service.Status() == types.Stopped {
		printer.PrintColored("Service "+service.Name()+" is already stopped.", printer.Service(), types.Red)
		return
	}
	service.SetStatus(types.Stopping)
	data := service.ServiceDataField
	data.(interfaces.Data).Stop(service, printer, clusterManager)
	service.SetStatus(types.Stopped)
}

func (service *Service) Start(printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	if (service.ImageNameField == "" || service.CurrentVersionField == "") && (service.ServiceType() != types.Harbor && service.ServiceType() != types.MongoDB) {
		printer.PrintColored("Service "+service.Name()+" is not configured correctly.", printer.Service(), types.Red)
		return
	}
	service.SetStatus(types.Starting)
	data := service.ServiceDataField
	data.(interfaces.Data).Start(service, printer, clusterManager)
	service.SetStatus(types.Running)
}

func (service *Service) Deploy(printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	if (service.ImageNameField == "" || service.CurrentVersionField == "") && (service.ServiceType() != types.Harbor && service.ServiceType() != types.MongoDB) {
		printer.PrintColored("Service "+service.Name()+" is not configured correctly.", printer.Service(), types.Red)
		return
	}
	service.SetStatus(types.Starting)
	data := service.ServiceDataField
	data.(interfaces.Data).Deploy(service, printer, clusterManager)
	service.updatePodInstances(clusterManager)
	service.SetStatus(types.Running)
}

func (service *Service) updatePodInstances(clusterManager interfaces.ClusterManager) {
	pods := clusterManager.GetPods("app="+service.Name(), service.Namespace())
	if len(pods) != 0 {
		for _, pod := range pods {
			service.AddPodInstance(pod.Name, &pod, string(pod.Status.Phase))
		}
	}
}
