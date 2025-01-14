package model

import (
	"NetManager/external/types"
	corev1 "k8s.io/api/core/v1"
)

type PodInstance struct {
	name string
	pod *corev1.Pod
	status string
}

func NewPodInstance(name string, pod *corev1.Pod, status string) *PodInstance {
	return &PodInstance{
		name: name,
		pod: pod,
		status: status,
	}
}

func CreateNewPodInstance(name string, pod *corev1.Pod) *PodInstance {
	return &PodInstance{
		name: name,
		pod: pod,
		status: types.Starting,
	}
}

func (instance PodInstance) Name() string {
	return instance.name
}

func (instance PodInstance) Pod() *corev1.Pod {
	return instance.pod
}

func (instance PodInstance) Status() string {
	return instance.status
}

