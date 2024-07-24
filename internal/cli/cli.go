package cli

import (
	"NetManager/internal/cli/manager"
	"NetManager/internal/service"
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type Console struct {
	service           service.Service
	inputBuffer       strings.Builder
	isRunning         bool
	isWaitingForInput bool
	printChan         chan string
	mutex             sync.Mutex
	wg                sync.WaitGroup
	commandManager    manager.CommandManager
}

func NewDefaultConsole() *Console {
	return &Console{
		service: service.Service{
			Name: "NetManager",
		},
		isRunning: true,
		printChan: make(chan string),
	}
}

func (console *Console) Init() *Console {
	console.commandManager = *manager.NewCommandManager(console)
	console.commandManager.Init() // Ensure commands are registered
	return console
}

func (console *Console) Run() {
	console.wg.Add(2)
	go console.handleInput()
	go console.printHandler()

	console.wg.Wait()
}

func (console *Console) Close() {
	console.isWaitingForInput = false
	console.isRunning = false
	close(console.printChan)
}

func (console *Console) Service() service.Service {
	return console.service
}

func (console *Console) handleInput() {
	defer console.wg.Done()
	reader := bufio.NewReader(os.Stdin)
	console.isWaitingForInput = true
	fmt.Print("> ")
	for console.isRunning {
		if !console.isWaitingForInput {
			continue
		}
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			console.clearLastInputLine()
			console.redrawInputLine()
			continue
		}
		console.isWaitingForInput = false
		go func(command string) {
			console.commandManager.ExecuteCommand(command)
		}(input)
	}
}

func (console *Console) printHandler() {
	defer console.wg.Done()
	for msg := range console.printChan {
		console.mutex.Lock()
		if console.isWaitingForInput {
			console.clearLastInputLine()
		}
		fmt.Println(msg)
		if !console.isRunning {
			break
		}
		console.redrawInputLine()
		console.isWaitingForInput = true
		console.mutex.Unlock()
	}
}

func (console *Console) Print(message string, service service.Service) {
	console.printChan <- fmt.Sprintf("[%s | %s]: %s", time.Now().Format("15:04:05"), service.Name, message)
}

func (console *Console) PrintColored(message string, service service.Service, color int) {
	coloredMessage := fmt.Sprintf("\033[38;5;%dm[%s | %s]: %s\033[0m", color, time.Now().Format("15:04:05"), service.Name, message)
	console.printChan <- coloredMessage
}

func (console *Console) clearLastInputLine() {
	fmt.Print("\033[1A\033[K")
}

func (console *Console) redrawInputLine() {
	fmt.Print("> ")
	fmt.Print(console.inputBuffer.String())
}
