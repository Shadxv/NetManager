package manager

import (
	"NetManager/external/cli"
	iDocker "NetManager/external/docker"
	"NetManager/external/service"
	docker "NetManager/internal/docker"
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/pkg/archive"
	"io"
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

func (manager *ImageManager) Client() iDocker.Client {
	return manager.client
}

func (manager *ImageManager) BuildImage(serviceModel service.Model) (bool, string) {
	serviceName := serviceModel.Name()

	manager.printer.Print("Starting building image for "+serviceName, manager.client.Service)

	contextDir := path.Join("services", "templates")
	tar, err := archive.Tar(contextDir, archive.Uncompressed)
	if err != nil {
		manager.printer.PrintColored("Error occured during building image. Operation canceled.", manager.client.Service, cli.Red)
		manager.printer.PrintColored(err.Error(), manager.client.Service, cli.Red)
		return false, ""
	}

	var newVersion string
	if len(serviceModel.AvailableVersions()) != 0 {
		newVersion, err = incrementVersion(serviceModel.AvailableVersions()[len(serviceModel.AvailableVersions())-1])
	} else {
		newVersion = "v0.1"
	}

	if err != nil {
		manager.printer.PrintColored(err.Error(), manager.client.Service, cli.Red)
		return false, ""
	}

	tag := serviceName + ":" + newVersion

	buildOptions := types.ImageBuildOptions{
		Tags:       []string{tag},
		Dockerfile: path.Join(serviceName, "Dockerfile"),
		NoCache:    true,
		Remove:     true,
	}

	build, err := manager.client.DockerClient.ImageBuild(context.Background(), tar, buildOptions)
	if err != nil {
		manager.printer.PrintColored("Error occured during building image. Operation canceled.", manager.client.Service, cli.Red)
		manager.printer.PrintColored(err.Error(), manager.client.Service, cli.Red)
		return false, ""
	}
	serviceModel.AddVersion(newVersion)

	imageID, err := getImageId(build.Body)

	if imageID == "" {
		manager.printer.PrintColored("Error occured during reading image id. Operation canceled.", manager.client.Service, cli.Red)
	}

	save, err := manager.client.DockerClient.ImageSave(context.Background(), []string{imageID})
	if err != nil {
		manager.printer.PrintColored("Error occured during saving image. Operation canceled.", manager.client.Service, cli.Red)
		manager.printer.PrintColored(err.Error(), manager.client.Service, cli.Red)
		return false, ""
	}
	defer build.Body.Close()
	defer save.Close()
	scanner := bufio.NewScanner(save)
	for scanner.Scan() {
		scanner.Text()
	}

	manager.printer.Print("Finished building image for "+serviceName, manager.client.Service)
	return true, imageID
}

func (manager *ImageManager) TagImage(serviceModel service.Model, serviceManager service.Manager) (bool, iDocker.HarborData) {
	manager.printer.Print("Creating tag new image version...", manager.client.Service)

	harborService := serviceManager.GetService("harbor")
	if harborService == nil {
		manager.printer.PrintColored("Could not tag new image. Harbor service is not loaded.", manager.client.Service, cli.Red)
		return false, nil
	}

	harborData := iDocker.GetHarborData(harborService)
	if harborData == nil {
		manager.printer.PrintColored("Could not tag new image. Harbor data is corrupted.", manager.client.Service, cli.Red)
		return false, nil
	}
	version := serviceModel.AvailableVersions()[len(serviceModel.AvailableVersions())-1]
	targetTag := harborData.Domain() + "/" + harborData.ProjectName() + "/" + serviceModel.ImageName() + ":" + version
	err := manager.client.DockerClient.ImageTag(context.Background(), serviceModel.ImageName()+":"+version, targetTag)
	if err != nil {
		manager.printer.PrintColored("Error occured during setting new image tag.", manager.client.Service, cli.Red)
		manager.printer.PrintColored(err.Error(), manager.client.Service, cli.Red)
		return false, nil
	}
	manager.printer.Print("Successfully tagged new image version: "+targetTag, manager.client.Service)
	return true, harborData
}

func (manager *ImageManager) PushImage(serviceModel service.Model, harborData iDocker.HarborData) bool {
	manager.printer.Print("Pushing image to remote registry...", manager.client.Service)

	version := serviceModel.AvailableVersions()[len(serviceModel.AvailableVersions())-1]
	tag := serviceModel.ImageName() + ":" + version
	targetTag := harborData.Domain() + "/" + harborData.ProjectName() + "/" + tag

	auth := registry.AuthConfig{
		Username: harborData.Username(),
		Password: harborData.Password(),
	}
	encodedAuth, err := registry.EncodeAuthConfig(auth)

	if err != nil {
		manager.printer.PrintColored("Error occured during encoding auth config.", manager.client.Service, cli.Red)
		manager.printer.PrintColored(err.Error(), manager.client.Service, cli.Red)
		return false
	}

	push, err := manager.client.DockerClient.ImagePush(context.Background(), targetTag, image.PushOptions{RegistryAuth: encodedAuth})
	if err != nil {
		manager.printer.PrintColored("Error occured during pushing image.", manager.client.Service, cli.Red)
		manager.printer.PrintColored(err.Error(), manager.client.Service, cli.Red)
		return false
	}
	defer push.Close()
	scanner := bufio.NewScanner(push)
	for scanner.Scan() {
		scanner.Text()
	}

	manager.printer.Print("Successfully pushed image to remote registry.", manager.client.Service)
	return true
}

func (manager *ImageManager) RemoveImage(id string) {
	_, err := manager.client.DockerClient.ImageRemove(context.Background(), id, image.RemoveOptions{Force: true})
	if err != nil {
		manager.printer.PrintColored("Error occured during clearing images.", manager.client.Service, cli.Red)
		manager.printer.PrintColored(err.Error(), manager.client.Service, cli.Red)
		return
	}

	_, err = manager.client.DockerClient.ImagesPrune(context.Background(), filters.NewArgs(filters.KeyValuePair{Key: "dangling", Value: "true"}))
	if err != nil {
		manager.printer.PrintColored("Error occured during clearing dangling images.", manager.client.Service, cli.Red)
		manager.printer.PrintColored(err.Error(), manager.client.Service, cli.Red)
		return
	}
}

func (manager *ImageManager) FullDeployImage(serviceModel service.Model, serviceManager service.Manager) bool {
	success, id := manager.BuildImage(serviceModel)
	if !success {
		return false
	}

	success, data := manager.TagImage(serviceModel, serviceManager)
	if !success {
		manager.RemoveImage(id)
		return false
	}

	if !manager.PushImage(serviceModel, data) {
		manager.RemoveImage(id)
		return false
	}

	manager.RemoveImage(id)

	manager.printer.Print("Image build has finished successfully.", manager.client.Service)
	return true
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

type buildResponse struct {
	Stream string `json:"stream,omitempty"`
	Aux    struct {
		ID string `json:"ID,omitempty"`
	} `json:"aux,omitempty"`
}

func getImageId(body io.ReadCloser) (string, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return "", err
	}

	responseString := string(data)

	var imageID string
	reader := strings.NewReader(responseString)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		var resp buildResponse
		err := json.Unmarshal([]byte(line), &resp)
		if err != nil {
			continue
		}

		if resp.Aux.ID != "" {
			imageID = resp.Aux.ID
			break
		}
	}

	if imageID == "" {
		return "", fmt.Errorf("image ID not found")
	}

	return imageID, nil
}
