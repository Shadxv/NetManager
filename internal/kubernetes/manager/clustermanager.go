package manager

import (
	"NetManager/internal/cli/handler"
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"context"
	"fmt"
	"io"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

type ClusterManager struct {
	service       interfaces.Service
	printer       interfaces.Printer
	configManager interfaces.ConfigManager
	clientset     *kubernetes.Clientset
	namespace     *corev1.Namespace
}

func NetClusterManager(service interfaces.Service, printer interfaces.Printer, configManager interfaces.ConfigManager, clientset *kubernetes.Clientset) *ClusterManager {
	return &ClusterManager{
		service:       service,
		printer:       printer,
		configManager: configManager,
		clientset:     clientset,
	}
}

func (manager *ClusterManager) Init() bool {
	manager.getNamespace()
	if manager.namespace == nil {
		return false
	}

	return true
}

func (manager *ClusterManager) getNamespace() {
	namespaceName := manager.configManager.GetMainConfig().GetName()
	namespace, err := manager.clientset.CoreV1().Namespaces().Get(context.Background(), namespaceName, metav1.GetOptions{})

	if err != nil {
		manager.printer.Print("Creating new namespace...", manager.service)
		newNamespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespaceName,
			},
		}
		namespace, err = manager.clientset.CoreV1().Namespaces().Create(context.Background(), newNamespace, metav1.CreateOptions{})
		if handler.HandleError(manager.printer, "Error occured during loading config.", err, manager.service, true) {
			return
		}
		manager.printer.Print("Created new namespace: "+namespaceName, manager.service)
	}

	manager.namespace = namespace
}

func (manager *ClusterManager) GetDefaultNamespace() string {
	return manager.namespace.Name
}

func (manager *ClusterManager) CreateNamespace(name string) error {
	_, err := manager.clientset.CoreV1().Namespaces().Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		newNamespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: name,
			},
		}
		_, err = manager.clientset.CoreV1().Namespaces().Create(context.Background(), newNamespace, metav1.CreateOptions{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (manager *ClusterManager) CreateDeployment(deployment *appsv1.Deployment, namespace string) *appsv1.Deployment {
	createdDeployment, err := manager.clientset.AppsV1().Deployments(namespace).Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return nil
	}
	return createdDeployment
}

func (manager *ClusterManager) UpdateDeployment(deployment *appsv1.Deployment, namespace string) {
	_, err := manager.clientset.AppsV1().Deployments(namespace).Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return
	}
}

func (manager *ClusterManager) DeleteDeployment(name string, namespace string) {
	err := manager.clientset.AppsV1().Deployments(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return
	}
}

func (manager *ClusterManager) GetDeployment(name string, namespace string) *appsv1.Deployment {
	deployment, err := manager.GetDeploymentOrErr(name, namespace)
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return nil
	}
	return deployment
}

func (manager *ClusterManager) GetDeploymentOrErr(name string, namespace string) (*appsv1.Deployment, error) {
	deployment, err := manager.clientset.AppsV1().Deployments(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return deployment, nil
}

func (manager *ClusterManager) CreateStatefulSet(statefulset *appsv1.StatefulSet, namespace string) *appsv1.StatefulSet {
	createdStatefulSet, err := manager.clientset.AppsV1().StatefulSets(namespace).Create(context.Background(), statefulset, metav1.CreateOptions{})
	if err != nil {
		manager.printer.PrintColored(fmt.Sprintf(err.Error()+" | namespace: %s", namespace), manager.service, types.Red)
		return nil
	}
	return createdStatefulSet
}

func (manager *ClusterManager) UpdateStatefulSet(statefulset *appsv1.StatefulSet, namespace string) {
	_, err := manager.clientset.AppsV1().StatefulSets(namespace).Update(context.Background(), statefulset, metav1.UpdateOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return
	}
}

func (manager *ClusterManager) DeleteStatefulSet(name string, namespace string) {
	err := manager.clientset.AppsV1().StatefulSets(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return
	}
}

func (manager *ClusterManager) GetStatefulSet(name string, namespace string) *appsv1.StatefulSet {
	statefulset, err := manager.GetStatefulSetOrErr(name, namespace)
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return nil
	}
	return statefulset
}

func (manager *ClusterManager) GetStatefulSetOrErr(name string, namespace string) (*appsv1.StatefulSet, error) {
	statefulset, err := manager.clientset.AppsV1().StatefulSets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return statefulset, err
}

func (manager *ClusterManager) CreateService(service *corev1.Service, namespace string) *corev1.Service {
	createdService, err := manager.clientset.CoreV1().Services(namespace).Create(context.Background(), service, metav1.CreateOptions{})
	if err != nil {
		manager.printer.PrintColored(fmt.Sprintf(err.Error()+" | namespace: %s", namespace), manager.service, types.Red)
		return nil
	}
	return createdService
}

func (manager *ClusterManager) UpdateService(service *corev1.Service, namespace string) {
	_, err := manager.clientset.CoreV1().Services(namespace).Update(context.Background(), service, metav1.UpdateOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return
	}
}

func (manager *ClusterManager) DeleteService(name string, namespace string) {
	err := manager.clientset.CoreV1().Services(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return
	}
}

func (manager *ClusterManager) GetService(name string, namespace string) *corev1.Service {
	service, err := manager.GetServiceOrErr(name, namespace)
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return nil
	}
	return service
}

func (manager *ClusterManager) GetServiceOrErr(name string, namespace string) (*corev1.Service, error) {
	service, err := manager.clientset.CoreV1().Services(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (manager *ClusterManager) CreateConfigMap(configmap *corev1.ConfigMap, namespace string) *corev1.ConfigMap {
	createdConfigMap, err := manager.clientset.CoreV1().ConfigMaps(namespace).Create(context.Background(), configmap, metav1.CreateOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return nil
	}
	return createdConfigMap
}

func (manager *ClusterManager) UpdateConfigMap(configmap *corev1.ConfigMap, namespace string) {
	_, err := manager.clientset.CoreV1().ConfigMaps(namespace).Update(context.Background(), configmap, metav1.UpdateOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return
	}
}

func (manager *ClusterManager) DeleteConfigMap(name string, namespace string) {
	err := manager.clientset.CoreV1().ConfigMaps(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return
	}
}

func (manager *ClusterManager) GetConfigMap(name string, namespace string) *corev1.ConfigMap {
	configmap, err := manager.GetConfigMapOrErr(name, namespace)
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return nil
	}
	return configmap
}

func (manager *ClusterManager) GetConfigMapOrErr(name string, namespace string) (*corev1.ConfigMap, error) {
	configmap, err := manager.clientset.CoreV1().ConfigMaps(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return configmap, nil
}

func (manager *ClusterManager) CreateSecret(secret *corev1.Secret, namespace string) *corev1.Secret {
	createdSecret, err := manager.clientset.CoreV1().Secrets(namespace).Create(context.Background(), secret, metav1.CreateOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		manager.printer.PrintColored("Could not create secret", manager.service, types.Red)
		return nil
	}
	return createdSecret
}

func (manager *ClusterManager) UpdateSecret(secret *corev1.Secret, namespace string) {
	_, err := manager.clientset.CoreV1().Secrets(namespace).Update(context.Background(), secret, metav1.UpdateOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return
	}
}

func (manager *ClusterManager) DeleteSecret(name string, namespace string) {
	err := manager.clientset.CoreV1().Secrets(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return
	}
}

func (manager *ClusterManager) GetSecret(name string, namespace string) *corev1.Secret {
	secret, err := manager.GetSecretOrErr(name, namespace)
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return nil
	}
	return secret
}

func (manager *ClusterManager) GetSecretOrErr(name string, namespace string) (*corev1.Secret, error) {
	secret, err := manager.clientset.CoreV1().Secrets(namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return secret, nil
}

func (manager *ClusterManager) GetPods(labelSelector string, namespace string) []corev1.Pod {
	pods, err := manager.clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})

	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return make([]corev1.Pod, 0)
	}

	return pods.Items
}

func (manager *ClusterManager) GetPodsOrErr(labelSelector string, namespace string) ([]corev1.Pod, error) {
	pods, err := manager.clientset.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})

	if err != nil {
		return make([]corev1.Pod, 0), err
	}

	return pods.Items, nil
}

func (manager *ClusterManager) GetPodLogs(podName string, namespace string) (io.ReadCloser, error) {
	return manager.clientset.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{
		Follow: true,
	}).Stream(context.Background())
}

func (manager *ClusterManager) WatchPods(labelSelector string, namespace string) (watch.Interface, error) {
	return manager.clientset.CoreV1().Pods(namespace).Watch(context.Background(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})
}

func (manager *ClusterManager) DeletePod(name string, namespace string) error {
	return manager.clientset.CoreV1().Pods(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
}

func (manager *ClusterManager) GetNodes() ([]corev1.Node, error) {
	nodes, err := manager.clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nodes.Items, nil
}
