package main

import (
	"NetManager/api"
	"NetManager/internal/cli"
	configManager "NetManager/internal/config/manager"
	dockerManager "NetManager/internal/docker/manager"
	harborModel "NetManager/internal/docker/model"
	"NetManager/internal/kubernetes"
	dataModel "NetManager/internal/minecraft/model"
	"NetManager/internal/module"
	"NetManager/internal/module/logger"
	mongoModel "NetManager/internal/mongodb/model"
	"NetManager/internal/redis"
	redisModel "NetManager/internal/redis/model"
	serviceManager "NetManager/internal/service/manager"
	serviceModel "NetManager/internal/service/model"
	"NetManager/pkg/types"
	"NetManager/pkg/util"
	"encoding/gob"
	"fmt"
	"sync"
)

func main() {
	var mainWaitGroup sync.WaitGroup

	err := util.CreateBaseDirStructure()
	if err != nil {
		fmt.Println("Fatal Error: " + err.Error())
		return
	}

	InitGob()
	log := logger.GetInstance()
	log.Init()

	moduleManager := module.NewModuleManager(&mainWaitGroup)
	defer moduleManager.Disable()

	var consoleModule = cli.NewDefaultConsole(&mainWaitGroup)
	moduleManager.AddModule(consoleModule)

	var configModule *configManager.ConfigManager
	if m, err := util.LoadData[configManager.ConfigManager](types.Config); err == nil {
		configModule = &m
	} else {
		configModule = configManager.NewConfigManager()
	}
	moduleManager.AddModule(configModule)

	var imagesModule = dockerManager.NewImageManager()
	moduleManager.AddModule(imagesModule)

	var clusterModule = kubernetes.NewClient()
	moduleManager.AddModule(clusterModule)

	var serviceModule *serviceManager.ServiceManager
	if m, err := util.LoadData[serviceManager.ServiceManager](types.Services); err == nil {
		serviceModule = &m
		for _, service := range serviceModule.Services {
			service.InitNonSavedFields()
		}
	} else {
		serviceModule = serviceManager.NewServiceManager()
	}
	moduleManager.AddModule(serviceModule)

	var redisModule = redis.NewRedisClient()
	moduleManager.AddModule(redisModule)

	var gatewayServer = api.NewServer(&mainWaitGroup)
	moduleManager.AddModule(gatewayServer)

	moduleManager.Init()

	mainWaitGroup.Wait()
}

func InitGob() {
	gob.Register(&dataModel.PaperData{})
	gob.Register(&dataModel.VelocityData{})
	gob.Register(&serviceModel.Service{})
	gob.Register(&harborModel.HarborData{})
	gob.Register(&redisModel.RedisData{})
	gob.Register(&mongoModel.MongoData{})
}
