package manager

import (
	"NetManager/external/cli"
	"NetManager/internal/docker"
	"bytes"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"os"
	"path"
	"path/filepath"
)

type ImageManager struct {
	printer cli.Printer
	Client  *docker.Client
}

func NewImageManager(printer cli.Printer) *ImageManager {
	return &ImageManager{
		printer: printer,
		Client:  docker.NewClient(printer),
	}
}

func (manager *ImageManager) Init() {
	manager.Client.Init()
}

func (manager *ImageManager) BuildImage(serviceName string) {
	manager.printer.Print("Starting building image for "+serviceName, manager.Client.Service)

	dockerfile := manager.buildDockerFile(serviceName)
	if dockerfile == "" {
		return
	}

	manager.printer.Print("\n"+dockerfile, manager.Client.Service)

	reader := bytes.NewReader([]byte(dockerfile))

	contextDir := path.Join("services", "templates", serviceName)
	contextFile, err := os.Open(contextDir)
	if err != nil {
		manager.printer.PrintColored("Unable to open context dir.", manager.Client.Service, cli.Red)
		return
	}
	defer contextFile.Close()

	buildOptions := types.ImageBuildOptions{
		Tags: []string{"dreammc-server:latest"},
	}

	_, err = manager.Client.DockerClient.ImageBuild(context.Background(), reader, buildOptions)
	if err != nil {
		manager.printer.PrintColored("Error occured during building image. Operation canceled.", manager.Client.Service, cli.Red)
		manager.printer.PrintColored(err.Error(), manager.Client.Service, cli.Red)
		return
	}
	manager.printer.Print("Finished building image for "+serviceName, manager.Client.Service)
}

func (manager *ImageManager) buildDockerFile(serviceName string) string {

	contextpath := filepath.Join("services", "templates", serviceName)

	return fmt.Sprintf(`
FROM eclipse-temurin:21-jre
WORKDIR /dreammc
COPY %s/ /dreammc/
RUN echo "eula=true" > eula.txt
CMD ["java -Xmx2G -Xms1G -jar server.jar nogui"]
`, contextpath)
}
