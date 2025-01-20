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
	GetDeploymentOrErr(name string) (*appsv1.Deployment, error)

	CreateStatefulSet(statefulSet *appsv1.StatefulSet) *appsv1.StatefulSet
	UpdateStatefulSet(statefulSet *appsv1.StatefulSet)
	DeleteStatefulSet(name string)
	GetStatefulSet(name string) *appsv1.StatefulSet
	GetStatefulSetOrErr(name string) (*appsv1.StatefulSet, error)

	CreateService(service *corev1.Service) *corev1.Service
	UpdateService(service *corev1.Service)
	DeleteService(name string)
	GetService(name string) *corev1.Service
	GetServiceOrErr(name string) (*corev1.Service, error)

	CreateConfigMap(configmap *corev1.ConfigMap) *corev1.ConfigMap
	UpdateConfigMap(configmap *corev1.ConfigMap)
	DeleteConfigMap(name string)
	GetConfigMap(name string) *corev1.ConfigMap
	GetConfigMapOrErr(name string) (*corev1.ConfigMap, error)

	CreateSecret(secret *corev1.Secret) *corev1.Secret
	UpdateSecret(secret *corev1.Secret)
	DeleteSecret(name string)
	GetSecret(name string) *corev1.Secret
	GetSecretOrErr(name string) (*corev1.Secret, error)

	GetPods(labelSelector string) []corev1.Pod
}
