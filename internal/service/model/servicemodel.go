package model

import (
	"NetManager/internal/kubernetes/model"
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	corev1 "k8s.io/api/core/v1"
	"slices"
)

type Service struct {
	name              string
	serviceType       string
	status            string
	imageName         string
	currentVersion    string
	availableVersions []string
	appConfig         interface{}
	podInstances      map[string]interfaces.PodInstance
	netService        *corev1.Service
	serviceData       interface{}
}

func NewService(name string, serviceType string, status string, imageName string, currentVersion string, availableVersions []string, serviceData interface{}) *Service {
	return &Service{
		name:              name,
		serviceType:       serviceType,
		status:            status,
		imageName:         imageName,
		currentVersion:    currentVersion,
		availableVersions: availableVersions,
		podInstances:      make(map[string]interfaces.PodInstance),
		serviceData:       serviceData,
	}
}

func CreateNewService(name string, serviceType string, serviceData interface{}) *Service {
	return &Service{
		name:         name,
		serviceType:  serviceType,
		status:       types.Stopped,
		imageName:    name,
		podInstances: make(map[string]interfaces.PodInstance),
		serviceData:  serviceData,
	}
}

func (service *Service) Name() string {
	return service.name
}

func (service *Service) ServiceType() string {
	return service.serviceType
}

func (service *Service) Status() string {
	return service.status
}

func (service *Service) SetStatus(status string) {
	service.status = status
}

func (service *Service) ImageName() string {
	return service.imageName
}

func (service *Service) CurrentVersion() string {
	return service.currentVersion
}

func (service *Service) SetCurrentVersion(currentVersion string) bool {
	if slices.Contains(service.availableVersions, currentVersion) {
		service.currentVersion = currentVersion
		return true
	}
	return false
}

func (service *Service) AvailableVersions() []string {
	versionsCopy := make([]string, len(service.availableVersions))
	copy(versionsCopy, service.availableVersions)
	return versionsCopy
}

func (service *Service) AddVersion(version string) {
	service.availableVersions = append(service.availableVersions, version)
}

func (service *Service) AppConfig() interface{} {
	return service.appConfig
}

func (service *Service) SetAppConfig(config interface{}) {
	service.appConfig = config
}

func (service *Service) RemoveAppConfig() {
	service.appConfig = nil
}

func (service *Service) PodInstances() map[string]interfaces.PodInstance {
	return service.podInstances
}

func (service *Service) addPodInstance(instance interfaces.PodInstance) {
	service.podInstances[instance.Name()] = instance
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
	delete(service.podInstances, name)
}

func (service *Service) NetService() *corev1.Service {
	return service.netService
}

func (service *Service) SetNetSevice(netService *corev1.Service) {
	service.netService = netService
}

func (service *Service) RemoveNetService() {
	service.netService = nil
}

func (service *Service) ServiceData() interface{} {
	return service.serviceData
}

func (service *Service) Build(printer interfaces.Printer) {
	data := service.serviceData
	data.(interfaces.Data).Build(service, printer)
}

func (service *Service) Update(printer interfaces.Printer) {
	data := service.serviceData
	data.(interfaces.Data).Update(service, printer)
}

func (service *Service) Stop(printer interfaces.Printer) {
	data := service.serviceData
	data.(interfaces.Data).Stop(service, printer)
}

func (service *Service) Start(printer interfaces.Printer) {
	data := service.serviceData
	data.(interfaces.Data).Start(service, printer)
}

func (service *Service) Deploy(printer interfaces.Printer) {
	data := service.serviceData
	data.(interfaces.Data).Deploy(service, printer)
}
