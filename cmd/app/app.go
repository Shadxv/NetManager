package main

import (
	"NetManager/internal/cli"
	cmdManager "NetManager/internal/config/manager"
	"NetManager/internal/kubernetes"
	serviceManager "NetManager/internal/minecraft/manager"
	"sync"
)

var Console *cli.Console
var KubernetesClient *kubernetes.Client
var ConfigManager *cmdManager.ConfigManager
var ServiceManager *serviceManager.ServiceManager

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

	KubernetesClient = kubernetes.NewClient(Console, ConfigManager)
	KubernetesClient.Load()
	if !KubernetesClient.IsLoaded {
		Console.CloseGracefully("App is shutting down...")
		return
	}
	//KubernetesClient.DeployRedis()

	if &(Console.CommandManager) == nil {
		Console.CloseGracefully("App is shutting down...")
		return
	}
	ServiceManager = serviceManager.NewServiceManager(Console, Console.CommandManager)
	ServiceManager.Init()

	wg.Wait()
}
