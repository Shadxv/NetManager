package docker

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	dockerClient "github.com/docker/docker/client"
)

type Client struct {
	DockerClient *dockerClient.Client
	printer      interfaces.Printer
	Service      interfaces.Service
}

func NewClient(printer interfaces.Printer) *Client {
	client := Client{
		printer: printer,
		Service: interfaces.Service{
			Name: "Docker DockerClient",
		},
	}
	return &client
}

func (client *Client) Init() {
	apiClient, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv)

	if err != nil {
		client.printer.PrintColored(err.Error(), client.Service, types.Red)
		client.printer.CloseGracefully("App is shutting down...")
		return
	}

	client.DockerClient = apiClient
}

func (client *Client) Close() {
	client.DockerClient.Close()
}
