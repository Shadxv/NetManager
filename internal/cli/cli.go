package cli

import (
	"NetManager/internal/cli/manager"
	"NetManager/internal/service"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
)

type Console struct {
	mainWaitGroup     *sync.WaitGroup
	service           service.Service
	isRunning         bool
	isWaitingForInput bool
	isClosing         bool
	printChan         chan string
	mutex             sync.Mutex
	wg                sync.WaitGroup
	commandManager    manager.CommandManager
	inputBuffer       string
	cursorPos         int
}

func NewDefaultConsole(mainWaitGroup *sync.WaitGroup) *Console {
	return &Console{
		mainWaitGroup: mainWaitGroup,
		service: service.Service{
			Name: "NetManager",
		},
		isRunning:   true,
		printChan:   make(chan string),
		inputBuffer: "",
		cursorPos:   0,
	}
}

func (console *Console) Init() *Console {
	console.commandManager = *manager.NewCommandManager(console)
	console.commandManager.Init()
	return console
}

func (console *Console) Run() {
	console.wg.Add(3)
	go console.handleInput()
	go console.printHandler()

	console.wg.Wait()
}

func (console *Console) Close() {
	console.mutex.Lock()
	defer console.mutex.Unlock()

	if !console.isRunning {
		return
	}

	console.isWaitingForInput = false
	console.isRunning = false
	close(console.printChan)
	keyboard.Close()
	console.mainWaitGroup.Done()
}

func (console *Console) Service() service.Service {
	return console.service
}

func (console *Console) handleInput() {
	defer console.wg.Done()
	if err := keyboard.Open(); err != nil {
		panic(err)
	}

	console.isWaitingForInput = true
	console.printPrompt()

	for console.isRunning {
		char, key, err := keyboard.GetKey()
		if !console.isRunning {
			return
		}

		if err != nil {
			panic(err)
		}

		switch key {
		case keyboard.KeyEnter:
			console.executeCommand()
		case keyboard.KeyBackspace, keyboard.KeyBackspace2:
			console.handleBackspace()
		case keyboard.KeyCtrlC:
			console.Close()
			return
		default:
			if char != 0 {
				console.handleChar(char)
			}
		}
	}
}

func (console *Console) executeCommand() {
	command := strings.TrimSpace(console.inputBuffer)
	console.inputBuffer = ""
	console.cursorPos = 0
	fmt.Println()
	if command != "" {
		go console.commandManager.ExecuteCommand(command)
	}
	console.printPrompt()

}

func (console *Console) handleBackspace() {
	if console.cursorPos > 0 {
		console.inputBuffer = console.inputBuffer[:console.cursorPos-1] + console.inputBuffer[console.cursorPos:]
		console.cursorPos--
		console.redrawInputLine()
	}
}

func (console *Console) handleChar(char rune) {
	console.inputBuffer = console.inputBuffer[:console.cursorPos] + string(char) + console.inputBuffer[console.cursorPos:]
	console.cursorPos++
	console.redrawInputLine()
}

func (console *Console) redrawInputLine() {
	fmt.Printf("\r> %s", console.inputBuffer)
	fmt.Printf("\033[%dG", console.cursorPos+3)
}

func (console *Console) printHandler() {
	defer console.wg.Done()
	for msg := range console.printChan {
		console.mutex.Lock()
		if console.isWaitingForInput {
			fmt.Print("\r\033[K")
			fmt.Println(msg)
			if !console.isClosing {
				console.redrawInputLine()
			}
		} else {
			fmt.Print("\r\033[K")
			fmt.Println(msg)
		}
		console.mutex.Unlock()
	}
}

func (console *Console) printPrompt() {
	fmt.Print("> ")
}

func (console *Console) Print(message string, service service.Service) {
	console.printChan <- fmt.Sprintf("[%s | %s]: %s", time.Now().Format("15:04:05"), service.Name, message)
}

func (console *Console) PrintColored(message string, service service.Service, color int) {
	coloredMessage := fmt.Sprintf("\033[38;5;%dm[%s | %s]: %s\033[0m", color, time.Now().Format("15:04:05"), service.Name, message)
	console.printChan <- coloredMessage
}

func (console *Console) SetClosingStatus() {
	console.isClosing = true
}
