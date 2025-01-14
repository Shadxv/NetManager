package main

import (
	"NetManager/internal/cli"
	cmdManager "NetManager/internal/config/manager"
	dockerManager "NetManager/internal/docker/manager"
	harborModel "NetManager/internal/docker/model"
	"NetManager/internal/kubernetes"
	serviceManager "NetManager/internal/service/manager"
	"sync"
)

var Console *cli.Console
var KubernetesClient *kubernetes.Client
var ConfigManager *cmdManager.ConfigManager
var ServiceManager *serviceManager.ServiceManager
var ImageManager *dockerManager.ImageManager

func main() {
	var wg sync.WaitGroup

	Console = cli.NewDefaultConsole(&wg)
	Console.Init()

	wg.Add(1)
	go func() {
		defer wg.Done()
		Console.Run()
	}()

	ConfigManager = cmdManager.NewConfigManager(Console)
	ConfigManager.Init()

	ImageManager = dockerManager.NewImageManager(Console)
	ImageManager.Init()
	defer ImageManager.Client.Close()

	KubernetesClient = kubernetes.NewClient(Console)
	KubernetesClient.Connect()
	KubernetesClient.Init(ConfigManager)
	if !KubernetesClient.IsLoaded() {
		Console.CloseGracefully("App is shutting down...")
		return
	}
	//KubernetesClient.DeployRedis()

	if &(Console.CommandManager) == nil {
		Console.CloseGracefully("App is shutting down...")
		return
	}
	ServiceManager = serviceManager.CreateNewServiceManager(Console)
	ServiceManager.Init(Console.CommandManager, ImageManager, KubernetesClient)

	harborModel.CreateHarborService(
		ConfigManager.GetHarborConfig(),
		ServiceManager,
		KubernetesClient.ClusterManager(),
	)

	wg.Wait()
}
