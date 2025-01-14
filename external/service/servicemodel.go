package service

import (
	"NetManager/external/kubernetes"
	corev1 "k8s.io/api/core/v1"
)

type Model interface {
	Name() string
	ServiceType() string
	Status() string
	SetStatus(status string)
	ImageName() string
	CurrentVersion() string
	SetCurrentVersion(currentVersion string) bool
	AvailableVersions() []string
	AddVersion(version string)
	AppConfig() *interface{}
	SetAppConfig(config *interface{})
	RemoveAppConfig()
	PodInstances() map[string]kubernetes.PodInstance
	AddPodInstance(name string, pod *corev1.Pod, status string) kubernetes.PodInstance
	CreatePodInstance(name string, pod *corev1.Pod) kubernetes.PodInstance
	RemovePodInstance(name string)
	NetService() *corev1.Service
	SetNetSevice(netService *corev1.Service)
	RemoveNetService()
	ServiceData() *interface{}
}
