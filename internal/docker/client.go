package docker

import (
	"NetManager/external/cli"
	dockerClient "github.com/docker/docker/client"
)

type Client struct {
	DockerClient *dockerClient.Client
	printer      cli.Printer
	Service      cli.Service
}

func NewClient(printer cli.Printer) *Client {
	client := Client{
		printer: printer,
		Service: cli.Service{
			Name: "Docker Client",
		},
	}
	return &client
}

func (client *Client) Init() {
	apiClient, err := dockerClient.NewClientWithOpts(dockerClient.FromEnv)

	if err != nil {
		client.printer.PrintColored(err.Error(), client.Service, cli.Red)
		client.printer.CloseGracefully("App is shutting down...")
		return
	}

	client.DockerClient = apiClient
}

func (client *Client) Close() {
	client.DockerClient.Close()
}
