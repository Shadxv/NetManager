package interfaces

import corev1 "k8s.io/api/core/v1"

type PodInstance interface {
	Name() string
	Pod() *corev1.Pod
	Status() string
	SetStatus(status string)
	InternalIP() string
	SetInternalIP(ip string)
	ExternalIP() string
	SetExternalIP(ip string)
	Port() int32
	SetPort(port int32)
	Logs() []string
	AddLog(line string)
	UpdatePod(pod *corev1.Pod)
}
