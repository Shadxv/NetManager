package kubernetes

import (
	"os"
	"path/filepath"

	"NetManager/internal/cli/model"
	"NetManager/internal/service"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Client struct {
	service   service.Service
	printer   model.Printer
	config    *rest.Config
	clientset *kubernetes.Clientset
	IsLoaded  bool
}

func NewClient(printer model.Printer) *Client {
	return &Client{
		service: service.Service{
			Name: "Kubernetes",
		},
		printer:  printer,
		IsLoaded: false,
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
		if err != nil {
			client.printer.PrintColored("Config not found. Check if Kubernetes Cluster is set up and properly configured!", client.service, model.Red)
			client.printer.PrintColored(err.Error(), client.service, model.Red)
			return
		}
	}
	client.config = config
	client.printer.Print("Loaded config", client.service)

	client.printer.Print("Connecting to Kubernetes Cluster...", client.service)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		client.printer.PrintColored("Could not connect to Kubernetes Cluster. Check if Kuberenetes is working!", client.service, model.Red)
		client.printer.PrintColored(err.Error(), client.service, model.Red)
		return
	}

	client.clientset = clientset
	client.IsLoaded = true
	client.printer.Print("Connected to Kubernetes Cluster", client.service)
}

// import (
// 	"context"
// 	"os"
// 	"path/filepath"

// 	"NetManager/internal/cli/model"

// 	corev1 "k8s.io/api/core/v1"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// 	"k8s.io/client-go/rest"
// 	"k8s.io/client-go/tools/clientcmd"
// )

// // NewKubernetesClient creates a new Kubernetes client
// func NewKubernetesClient(printer *model.Printer) (*[]corev1.Namespace, error) {
// 	// Try to load in-cluster config
// 	config, err := rest.InClusterConfig()
// 	if err != nil {
// 		// If that fails, try loading config from ~/.kube/config
// 		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
// 		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}

// 	// Create the clientset
// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// List all namespaces
// 	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})

// 	return &namespaces.Items, nil
// }
