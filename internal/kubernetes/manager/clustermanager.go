package manager

import (
	"NetManager/internal/cli/handler"
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"context"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (manager *ClusterManager) CreateDeployment(deployment *appsv1.Deployment) *appsv1.Deployment {
	createdDeployment, err := manager.clientset.AppsV1().Deployments(manager.namespace.GetNamespace()).Create(context.Background(), deployment, metav1.CreateOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return nil
	}
	return createdDeployment
}

func (manager *ClusterManager) UpdateDeployment(deployment *appsv1.Deployment) {
	_, err := manager.clientset.AppsV1().Deployments(manager.namespace.GetNamespace()).Update(context.Background(), deployment, metav1.UpdateOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return
	}
}

func (manager *ClusterManager) DeleteDeployment(name string) {
	err := manager.clientset.AppsV1().Deployments(manager.namespace.GetNamespace()).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return
	}
}

func (manager *ClusterManager) GetDeployment(name string) *appsv1.Deployment {
	deployment, err := manager.clientset.AppsV1().Deployments(manager.namespace.GetNamespace()).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return nil
	}
	return deployment
}

func (manager *ClusterManager) CreateStatefulSet(statefulset *appsv1.StatefulSet) *appsv1.StatefulSet {
	createdStatefulSet, err := manager.clientset.AppsV1().StatefulSets(manager.namespace.GetNamespace()).Create(context.Background(), statefulset, metav1.CreateOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return nil
	}
	return createdStatefulSet
}

func (manager *ClusterManager) UpdateStatefulSet(statefulset *appsv1.StatefulSet) {
	_, err := manager.clientset.AppsV1().StatefulSets(manager.namespace.GetNamespace()).Update(context.Background(), statefulset, metav1.UpdateOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return
	}
}

func (manager *ClusterManager) DeleteStatefulSet(name string) {
	err := manager.clientset.AppsV1().StatefulSets(manager.namespace.GetNamespace()).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return
	}
}

func (manager *ClusterManager) GetStatefulSet(name string) *appsv1.StatefulSet {
	statefulset, err := manager.clientset.AppsV1().StatefulSets(manager.namespace.GetNamespace()).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return nil
	}
	return statefulset
}

func (manager *ClusterManager) CreateService(service *corev1.Service) *corev1.Service {
	createdService, err := manager.clientset.CoreV1().Services(manager.namespace.GetNamespace()).Create(context.Background(), service, metav1.CreateOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return nil
	}
	return createdService
}

func (manager *ClusterManager) UpdateService(service *corev1.Service) {
	_, err := manager.clientset.CoreV1().Services(manager.namespace.GetNamespace()).Update(context.Background(), service, metav1.UpdateOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return
	}
}

func (manager *ClusterManager) DeleteService(name string) {
	err := manager.clientset.CoreV1().Services(manager.namespace.GetNamespace()).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return
	}
}

func (manager *ClusterManager) GetService(name string) *corev1.Service {
	service, err := manager.clientset.CoreV1().Services(manager.namespace.GetNamespace()).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return nil
	}
	return service
}

func (manager *ClusterManager) CreateConfigMap(service *corev1.ConfigMap) *corev1.ConfigMap {
	createdConfigMap, err := manager.clientset.CoreV1().ConfigMaps(manager.namespace.GetNamespace()).Create(context.Background(), service, metav1.CreateOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return nil
	}
	return createdConfigMap
}

func (manager *ClusterManager) UpdateConfigMap(service *corev1.ConfigMap) {
	_, err := manager.clientset.CoreV1().ConfigMaps(manager.namespace.GetNamespace()).Update(context.Background(), service, metav1.UpdateOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return
	}
}

func (manager *ClusterManager) DeleteConfigMap(name string) {
	err := manager.clientset.CoreV1().ConfigMaps(manager.namespace.GetNamespace()).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return
	}
}

func (manager *ClusterManager) GetConfigMap(name string) *corev1.ConfigMap {
	configmap, err := manager.clientset.CoreV1().ConfigMaps(manager.namespace.GetNamespace()).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return nil
	}
	return configmap
}

func (manager *ClusterManager) GetPods(labelSelector string) []corev1.Pod {
	pods, err := manager.clientset.CoreV1().Pods(manager.namespace.Namespace).List(context.Background(), metav1.ListOptions{
		LabelSelector: labelSelector,
	})

	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.service, types.Red)
		return make([]corev1.Pod, 0)
	}

	return pods.Items
}
