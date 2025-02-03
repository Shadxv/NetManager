package main

import (
	"NetManager/internal/cli"
	cmdManager "NetManager/internal/config/manager"
	dockerManager "NetManager/internal/docker/manager"
	harborModel "NetManager/internal/docker/model"
	"NetManager/internal/kubernetes"
	mongodbModel "NetManager/internal/mongodb/model"
	redisModel "NetManager/internal/redis/model"
	serviceManager "NetManager/internal/service/manager"
	"sync"
)

// TODO: Later add helm client to auto deploy charts like harbor and longhorn

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
	defer ImageManager.Client().Close()

	KubernetesClient = kubernetes.NewClient(Console)
	KubernetesClient.Connect()
	KubernetesClient.Init(ConfigManager)
	if !KubernetesClient.IsLoaded() {
		Console.CloseGracefully("App is shutting down...")
		return
	}

	if &(Console.CommandManager) == nil {
		Console.CloseGracefully("App is shutting down...")
		return
	}
	ServiceManager = serviceManager.CreateNewServiceManager(Console)
	ServiceManager.Init(Console.CommandManager, ImageManager, KubernetesClient)

	harborModel.CreateHarborService(
		Console,
		ConfigManager.GetHarborConfig(),
		ServiceManager,
		KubernetesClient.ClusterManager(),
	)

	redisModel.CreateNewRedisService(
		Console,
		ConfigManager.GetRedisConfig(),
		ServiceManager,
		KubernetesClient.ClusterManager(),
	)

	mongodbModel.CreateNewMongoService(
		ConfigManager.GetMongoConfig(),
		ServiceManager,
	)

	wg.Wait()
}
