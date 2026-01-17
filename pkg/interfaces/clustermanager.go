package interfaces

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

type ClusterManager interface {
	GetDefaultNamespace() string

	CreateDeployment(deployment *appsv1.Deployment, namespace string) *appsv1.Deployment
	UpdateDeployment(deployment *appsv1.Deployment, namespace string)
	DeleteDeployment(name string, namespace string)
	GetDeployment(name string, namespace string) *appsv1.Deployment
	GetDeploymentOrErr(name string, namespace string) (*appsv1.Deployment, error)

	CreateStatefulSet(statefulSet *appsv1.StatefulSet, namespace string) *appsv1.StatefulSet
	UpdateStatefulSet(statefulSet *appsv1.StatefulSet, namespace string)
	DeleteStatefulSet(name string, namespace string)
	GetStatefulSet(name string, namespace string) *appsv1.StatefulSet
	GetStatefulSetOrErr(name string, namespace string) (*appsv1.StatefulSet, error)

	CreateService(service *corev1.Service, namespace string) *corev1.Service
	UpdateService(service *corev1.Service, namespace string)
	DeleteService(name string, namespace string)
	GetService(name string, namespace string) *corev1.Service
	GetServiceOrErr(name string, namespace string) (*corev1.Service, error)

	CreateConfigMap(configmap *corev1.ConfigMap, namespace string) *corev1.ConfigMap
	UpdateConfigMap(configmap *corev1.ConfigMap, namespace string)
	DeleteConfigMap(name string, namespace string)
	GetConfigMap(name string, namespace string) *corev1.ConfigMap
	GetConfigMapOrErr(name string, namespace string) (*corev1.ConfigMap, error)

	CreateSecret(secret *corev1.Secret, namespace string) *corev1.Secret
	UpdateSecret(secret *corev1.Secret, namespace string)
	DeleteSecret(name string, namespace string)
	GetSecret(name string, namespace string) *corev1.Secret
	GetSecretOrErr(name string, namespace string) (*corev1.Secret, error)

	GetPods(labelSelector string, namespace string) []corev1.Pod
	GetPodsOrErr(labelSelector string, namespace string) ([]corev1.Pod, error)

	GetNodes() ([]corev1.Node, error)
}
