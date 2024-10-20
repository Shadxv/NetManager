package kubernetes

import (
	"context"
	"os"
	"path/filepath"

	"NetManager/internal/cli/handler"
	"NetManager/internal/cli/model"
	"NetManager/internal/config/manager"
	"NetManager/internal/kubernetes/config"
	"NetManager/internal/service"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	service       service.Service
	printer       model.Printer
	configManager *manager.ConfigManager
	config        *rest.Config
	clientset     *kubernetes.Clientset
	namespace     *corev1.Namespace
	IsLoaded      bool
}

func NewClient(printer model.Printer, configManager *manager.ConfigManager) *Client {
	return &Client{
		service: service.Service{
			Name: "Kubernetes",
		},
		printer:       printer,
		configManager: configManager,
		IsLoaded:      false,
	}
}

func (client *Client) checkIfLoaded() bool {
	return client.IsLoaded
}

func (client *Client) Load() {
	client.printer.Print("Trying to load config...", client.service)
	config, err := rest.InClusterConfig()
	if err != nil {
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if handler.HandleError(client.printer, "Config not found. Check if Kubernetes Cluster is set up and properly configured!", err, client.service, true) {
			return
		}
	}
	client.config = config
	client.printer.Print("Loaded config", client.service)

	client.printer.Print("Connecting to Kubernetes Cluster...", client.service)
	clientset, err := kubernetes.NewForConfig(config)
	if handler.HandleError(client.printer, "Could not connect to Kubernetes Cluster. Check if Kuberenetes is working!", err, client.service, true) {
		return
	}

	client.clientset = clientset
	client.IsLoaded = true
	client.printer.Print("Connected to Kubernetes Cluster", client.service)

	client.GetNamespace()
}

func (client *Client) GetNamespace() {
	namespaceName := client.configManager.GetMainConfig().Name
	namespace, err := client.clientset.CoreV1().Namespaces().Get(context.Background(), namespaceName, metav1.GetOptions{})

	if err != nil {
		client.printer.Print("Creating new namespace...", client.service)
		newNamespace := &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespaceName,
			},
		}
		namespace, err = client.clientset.CoreV1().Namespaces().Create(context.Background(), newNamespace, metav1.CreateOptions{})
		if handler.HandleError(client.printer, "Error occured during loading config.", err, client.service, true) {
			return
		}
		client.printer.Print("Created new namespace: "+namespaceName, client.service)
	}

	client.namespace = namespace
}

func (client *Client) DeployRedis() {
	configMap, deployment := config.GenerateRedisConfig(client.configManager.GetRedisConfig())

	_, err := client.clientset.CoreV1().ConfigMaps(client.configManager.GetMainConfig().Name).Create(context.TODO(), configMap, metav1.CreateOptions{})
	if handler.HandleError(client.printer, "Error occured during creating redis config map.", err, client.service, false) {
		return
	}

	_, err = client.clientset.AppsV1().Deployments(client.configManager.GetMainConfig().Name).Create(context.TODO(), deployment, metav1.CreateOptions{})
	if handler.HandleError(client.printer, "Error occured during redis deployment.", err, client.service, false) {
		return
	}
}
