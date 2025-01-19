package interfaces

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type ClusterManager interface {
	CreateDeployment(deployment *appsv1.Deployment) *appsv1.Deployment
	UpdateDeployment(deployment *appsv1.Deployment)
	DeleteDeployment(name string)
	GetDeployment(name string) *appsv1.Deployment

	CreateStatefulSet(statefulSet *appsv1.StatefulSet) *appsv1.StatefulSet
	UpdateStatefulSet(statefulSet *appsv1.StatefulSet)
	DeleteStatefulSet(name string)
	GetStatefulSet(name string) *appsv1.StatefulSet

	CreateService(service *corev1.Service) *corev1.Service
	UpdateService(service *corev1.Service)
	DeleteService(name string)
	GetService(name string) *corev1.Service

	CreateConfigMap(configmap *corev1.ConfigMap) *corev1.ConfigMap
	UpdateConfigMap(configmap *corev1.ConfigMap)
	DeleteConfigMap(name string)
	GetConfigMap(name string) *corev1.ConfigMap

	GetPods(labelSelector string) []corev1.Pod
}
