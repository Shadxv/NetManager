package model

import (
	"NetManager/pkg/types"
	"sync"

	corev1 "k8s.io/api/core/v1"
)

type PodInstance struct {
	name       string
	pod        *corev1.Pod
	status     string
	internalIP string
	externalIP string
	port       int32
	logs       []string
	logMutex   sync.RWMutex
}

func NewPodInstance(name string, pod *corev1.Pod, status string) *PodInstance {
	return &PodInstance{
		name:   name,
		pod:    pod,
		status: status,
		logs:   make([]string, 0, 1000),
	}
}

func CreateNewPodInstance(name string, pod *corev1.Pod) *PodInstance {
	return &PodInstance{
		name:   name,
		pod:    pod,
		status: types.Starting,
		logs:   make([]string, 0, 1000),
	}
}

func (instance *PodInstance) Name() string {
	return instance.name
}

func (instance *PodInstance) Pod() *corev1.Pod {
	return instance.pod
}

func (instance *PodInstance) Status() string {
	return instance.status
}

func (instance *PodInstance) SetStatus(status string) {
	instance.status = status
}

func (instance *PodInstance) InternalIP() string {
	return instance.internalIP
}

func (instance *PodInstance) SetInternalIP(ip string) {
	instance.internalIP = ip
}

func (instance *PodInstance) ExternalIP() string {
	return instance.externalIP
}

func (instance *PodInstance) SetExternalIP(ip string) {
	instance.externalIP = ip
}

func (instance *PodInstance) Port() int32 {
	return instance.port
}

func (instance *PodInstance) SetPort(port int32) {
	instance.port = port
}

func (instance *PodInstance) Logs() []string {
	instance.logMutex.RLock()
	defer instance.logMutex.RUnlock()
	logsCopy := make([]string, len(instance.logs))
	copy(logsCopy, instance.logs)
	return logsCopy
}

func (instance *PodInstance) AddLog(line string) {
	instance.logMutex.Lock()
	defer instance.logMutex.Unlock()

	if len(instance.logs) >= 1000 {
		instance.logs = instance.logs[1:]
	}
	instance.logs = append(instance.logs, line)
}

func (instance *PodInstance) UpdatePod(pod *corev1.Pod) {
	instance.pod = pod
}
