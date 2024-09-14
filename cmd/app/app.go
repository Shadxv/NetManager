package main

import (
	"NetManager/internal/cli"
	"NetManager/internal/kubernetes"
	"sync"
)

var Console *cli.Console
var KubernetesClient *kubernetes.Client

func main() {
	var wg sync.WaitGroup

	Console = cli.NewDefaultConsole(&wg)
	Console.Init()

	wg.Add(1)
	go func() {
		defer wg.Done()
		Console.Run()
	}()

	KubernetesClient = kubernetes.NewClient(Console)
	KubernetesClient.Load()
	if !KubernetesClient.IsLoaded {
		Console.SetClosingStatus()
		Console.Print("App is shutting down...", Console.Service())
		Console.Close()
		return
	}

	wg.Wait()
}
