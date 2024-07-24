package main

import "NetManager/internal/cli"

var Console *cli.Console

func main() {
	Console = cli.NewDefaultConsole()
	Console.Init().Run()
}
