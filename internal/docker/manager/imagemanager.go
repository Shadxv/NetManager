package manager

import (
	"NetManager/external/cli"
	"NetManager/external/service"
	"NetManager/internal/docker"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/pkg/archive"
	"path"
	"strconv"
	"strings"
)

type ImageManager struct {
	printer cli.Printer
	client  *docker.Client
}

func NewImageManager(printer cli.Printer) *ImageManager {
	return &ImageManager{
		printer: printer,
		client:  docker.NewClient(printer),
	}
}

func (manager *ImageManager) Init() {
	manager.client.Init()
}

func (manager *ImageManager) BuildImage(serviceModel service.Model) {
	serviceName := serviceModel.Name()

	manager.printer.Print("Starting building image for "+serviceName, manager.client.Service)

	contextDir := path.Join("services", "templates")
	tar, err := archive.Tar(contextDir, archive.Uncompressed)
	if err != nil {
		manager.printer.PrintColored("Error occured during building image. Operation canceled.", manager.client.Service, cli.Red)
		manager.printer.PrintColored(err.Error(), manager.client.Service, cli.Red)
		return
	}

	var newVersion string
	if len(serviceModel.AvailableVersions()) != 0 {
		newVersion, err = incrementVersion(serviceModel.AvailableVersions()[len(serviceModel.AvailableVersions())-1])
	} else {
		newVersion = "v0.1"
	}

	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.client.Service, cli.Red)
		return
	}

	buildOptions := types.ImageBuildOptions{
		Tags:       []string{serviceName + ":" + newVersion},
		Dockerfile: path.Base(path.Join(contextDir, serviceName, "Dockerfile")),
		NoCache:    true,
	}

	resp, err := manager.client.DockerClient.ImageBuild(context.Background(), tar, buildOptions)
	if err != nil {
		manager.printer.PrintColored("Error occured during building image. Operation canceled.", manager.client.Service, cli.Red)
		manager.printer.PrintColored(err.Error(), manager.client.Service, cli.Red)
		return
	}
	defer resp.Body.Close()

	manager.printer.Print("Finished building image for "+serviceName, manager.client.Service)
}

func incrementVersion(version string) (string, error) {
	if !strings.HasPrefix(version, "v") {
		return "", fmt.Errorf("corrupted version name: %s", version)
	}
	version = version[1:]
	parts := strings.Split(version, ".")
	if len(parts) != 2 {
		return "", fmt.Errorf("corrupted version name: %s", version)
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", fmt.Errorf("corrupted version name: %s", version)
	}
	minor++

	newVersion := fmt.Sprintf("v%s.%d", parts[0], minor)

	return newVersion, nil
}
