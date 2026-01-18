package interfaces

import (
	corev1 "k8s.io/api/core/v1"
)

type ServiceModel interface {
	Name() string
	ServiceType() string
	Status() string
	SetStatus(status string)
	ImageName() string
	Namespace() string
	CurrentVersion() string
	SetCurrentVersion(currentVersion string) bool
	AvailableVersions() []string
	AddVersion(version string)
	AppConfig() interface{}
	SetAppConfig(config interface{})
	RemoveAppConfig()
	PodInstances() map[string]PodInstance
	AddPodInstance(name string, pod *corev1.Pod, status string) PodInstance
	CreatePodInstance(name string, pod *corev1.Pod) PodInstance
	RemovePodInstance(name string)
	NetService() *corev1.Service
	SetNetSevice(netService *corev1.Service)
	RemoveNetService()
	ServiceData() interface{}
	Build(printer Printer, imageManager ImageManager, serviceManager ServiceManager)
	Update(printer Printer, clusterManager ClusterManager)
	Stop(printer Printer, clusterManager ClusterManager)
	Start(printer Printer, clusterManager ClusterManager)
	Deploy(printer Printer, clusterManager ClusterManager)
	SetBroadcaster(broadcaster Broadcaster)
	InitNonSavedFields()
}
