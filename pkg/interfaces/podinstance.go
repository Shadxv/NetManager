package interfaces

import corev1 "k8s.io/api/core/v1"

type PodInstance interface {
	Name() string
	Pod() *corev1.Pod
	Status() string
}
