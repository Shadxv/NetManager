package kubernetes

import (
	"NetManager/internal/cli/handler"
	"NetManager/internal/kubernetes/manager"
	"NetManager/internal/module"
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	service        interfaces.Service
	printer        interfaces.Printer
	config         *rest.Config
	clientset      *kubernetes.Clientset
	clusterManager *manager.ClusterManager
	isLoaded       bool

	// Module fields
	status        string
	moduleManager *module.Manager
}

func NewClient() *Client {
	return &Client{
		service: interfaces.Service{
			Name: "Kubernetes",
		},
		isLoaded: false,
		status:   types.Starting,
	}
}

func (client *Client) Init(moduleManager *module.Manager) {
	client.moduleManager = moduleManager

	printer, err := module.GetTypedModule[interfaces.Printer](client.moduleManager, types.Console)
	if err != nil {
		client.SetStatus(types.Disabled)
		return
	}
	client.printer = printer

	configManager, err := module.GetTypedModule[interfaces.ConfigManager](client.moduleManager, types.Config)
	if err != nil {
		client.printer.PrintColored("Critical error: ConfigManager module not found", client.service, types.Red)
		client.printer.Print(err.Error(), client.service)
		client.SetStatus(types.Disabled)
		return
	}

	client.Connect()

	if client.clientset != nil {
		client.clusterManager = manager.NetClusterManager(client.service, client.printer, configManager, client.clientset)
		client.isLoaded = client.clusterManager.Init()
	}

	client.SetStatus(types.Enabled)
}

func (client *Client) Disable(shutdown bool) {
	client.SetStatus(types.Disabled)
}

func (client *Client) Reload() {
	client.SetStatus(types.Starting)
	client.Connect()
	if client.moduleManager != nil {
		client.Init(client.moduleManager)
	}
	client.SetStatus(types.Enabled)
}

func (client *Client) SaveData() error {
	return nil
}

func (client *Client) LoadData() {}

func (client *Client) SetStatus(newStatus string) {
	client.status = newStatus
}

func (client *Client) Status() string {
	return client.status
}

func (client *Client) Type() string {
	return types.Kubernetes
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

func (client *Client) ClusterManager() interfaces.ClusterManager {
	return client.clusterManager
}
