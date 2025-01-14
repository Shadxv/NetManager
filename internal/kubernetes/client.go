package kubernetes

import (
	"NetManager/external/cli"
	iConfig "NetManager/external/config"
	"NetManager/internal/kubernetes/manager"
	"os"
	"path/filepath"

	"NetManager/internal/cli/handler"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	service        cli.Service
	printer        cli.Printer
	config         *rest.Config
	clientset      *kubernetes.Clientset
	clusterManager *manager.ClusterManager
	isLoaded       bool
}

func NewClient(printer cli.Printer) *Client {
	return &Client{
		service: cli.Service{
			Name: "Kubernetes",
		},
		printer:  printer,
		isLoaded: false,
	}
}

func (client *Client) Init(configManager iConfig.Manager) {
	if client.clientset == nil {
		return
	}

	client.clusterManager = manager.NetClusterManager(client.service, client.printer, configManager, client.clientset)
	client.isLoaded = client.clusterManager.Init()
}

func (client *Client) IsLoaded() bool {
	return client.isLoaded
}

func (client *Client) Connect() {
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
	client.isLoaded = true
	client.printer.Print("Connected to Kubernetes Cluster", client.service)
}

func (client *Client) ClusterManager() *manager.ClusterManager {
	return client.clusterManager
}

//func (client *Client) DeployRedis() {
//	configMap, deployment := config.GenerateRedisConfig(client.configManager.GetRedisConfig())
//
//	_, err := client.clientset.CoreV1().ConfigMaps(client.configManager.GetMainConfig().Name).Create(context.TODO(), configMap, metav1.CreateOptions{})
//	if handler.HandleError(client.printer, "Error occured during creating redis config map.", err, client.service, false) {
//		return
//	}
//
//	_, err = client.clientset.AppsV1().Deployments(client.configManager.GetMainConfig().Name).Create(context.TODO(), deployment, metav1.CreateOptions{})
//	if handler.HandleError(client.printer, "Error occured during redis deployment.", err, client.service, false) {
//		return
//	}
//}
//
//func (client *Client) DeployPaperService(serviceName string) {
//	deployment := config.GeneratePaperDeployment(serviceName)
//	_, err := client.clientset.AppsV1().Deployments(client.configManager.GetMainConfig().Name).Create(context.TODO(), deployment, metav1.CreateOptions{})
//	if err != nil {
//		client.printer.PrintColored("Error occured during "+serviceName+" deployment:", client.service, cli.Red)
//		client.printer.PrintColored(err.Error(), client.service, cli.Red)
//	}
//}
