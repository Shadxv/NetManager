package main

import (
	"NetManager/internal/cli"
	configManager "NetManager/internal/config/manager"
	dockerManager "NetManager/internal/docker/manager"
	"NetManager/internal/kubernetes"
	dataModel "NetManager/internal/minecraft/model"
	"NetManager/internal/module"
	"NetManager/internal/module/logger"
	"NetManager/internal/redis"
	serviceManager "NetManager/internal/service/manager"
	serviceModel "NetManager/internal/service/model"
	"NetManager/pkg/types"
	"NetManager/pkg/util"
	"encoding/gob"
	"fmt"
	"sync"
)

var Console *cli.Console
var KubernetesClient *kubernetes.Client
var ConfigManager *configManager.ConfigManager
var ServiceManager *serviceManager.ServiceManager
var ImageManager *dockerManager.ImageManager
var RedisClient *redis.Client

func main() {
	var mainWaitGroup sync.WaitGroup

	err := util.CreateBaseDirStructure()
	if err != nil {
		fmt.Println("Fatal Error: " + err.Error())
		return
	}

	log := logger.GetInstance()
	log.Init()

	moduleManager := module.NewModuleManager(&mainWaitGroup)

	var configModule *configManager.ConfigManager
	if m, err := util.LoadData[configManager.ConfigManager](types.Config); err == nil {
		configModule = &m
	} else {
		configModule = configManager.NewConfigManager()
	}
	moduleManager.AddModule(configModule)

	mainWaitGroup.Wait()
}

func InitGob() {
	gob.Register(dataModel.PaperData{})
	gob.Register(dataModel.VelocityData{})
	gob.Register(serviceModel.Service{})
}

//func main() {
//	var wg sync.WaitGroup
//
//	Console = cli.NewDefaultConsole(&wg)
//	Console.Init()
//
//	wg.Add(1)
//	go func() {
//		defer wg.Done()
//		Console.Run()
//	}()
//
//	ConfigManager = cmdManager.NewConfigManager(Console)
//	ConfigManager.Init()
//
//	ImageManager = dockerManager.NewImageManager(Console)
//	ImageManager.Init()
//	defer ImageManager.Client().Close()
//
//	KubernetesClient = kubernetes.NewClient(Console)
//	KubernetesClient.Connect()
//	KubernetesClient.Init(ConfigManager)
//	if !KubernetesClient.IsLoaded() {
//		Console.CloseGracefully("App is shutting down...")
//		return
//	}
//
//	if Console.CommandManager() == nil {
//		Console.CloseGracefully("App is shutting down...")
//		return
//	}
//	ServiceManager = serviceManager.CreateNewServiceManager(Console, ConfigManager)
//	ServiceManager.Init(Console.CommandManager, ImageManager, KubernetesClient)
//
//	harborModel.CreateHarborService(
//		Console,
//		ConfigManager.GetHarborConfig(),
//		ServiceManager,
//		KubernetesClient.ClusterManager(),
//	)
//
//	redisModel.CreateNewRedisService(
//		Console,
//		ConfigManager.GetRedisConfig(),
//		ServiceManager,
//		KubernetesClient.ClusterManager(),
//	)
//
//	mongodbModel.CreateNewMongoService(
//		ConfigManager.GetMongoConfig(),
//		ServiceManager,
//	)
//
//	RedisClient = redis.NewRedisClient(Console, KubernetesClient.ClusterManager())
//	RedisClient.Init(ServiceManager.GetService("redis"))
//	defer RedisClient.Close()
//
//	wg.Wait()
//}
