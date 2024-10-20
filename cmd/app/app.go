package main

import (
	"NetManager/internal/cli"
	"NetManager/internal/config/manager"
	"NetManager/internal/kubernetes"
	"sync"
)

var Console *cli.Console
var KubernetesClient *kubernetes.Client
var ConfigManager *manager.ConfigManager

func main() {
	var wg sync.WaitGroup

	Console = cli.NewDefaultConsole(&wg)
	Console.Init()

	wg.Add(1)
	go func() {
		defer wg.Done()
		Console.Run()
	}()

	ConfigManager = manager.NewConfigManager(Console)
	ConfigManager.Init()

	KubernetesClient = kubernetes.NewClient(Console, ConfigManager)
	KubernetesClient.Load()
	if !KubernetesClient.IsLoaded {
		Console.CloseGracefully("App is shutting down...")
		return
	}
	KubernetesClient.DeployRedis()

	wg.Wait()
}
