package cli

import (
	"NetManager/internal/cli/handler"
	"NetManager/internal/cli/manager"
	"NetManager/internal/cli/model"
	"NetManager/internal/service"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/eiannone/keyboard"
)

type Console struct {
	CommandManager *manager.CommandManager
	// Prevents app from closing befor it should (console is async)
	mainWaitGroup     *sync.WaitGroup
	service           service.Service
	isRunning         bool
	isWaitingForInput bool
	isClosing         bool
	isPaused          bool
	pendingMessages   []string
	printChan         chan string
	mutex             sync.RWMutex
	// To remove, I think it is useless
	consoleWaitGroup sync.WaitGroup
	// Allows printer to print last messages before closing app
	printHandlerWG sync.WaitGroup
	inputBuffer    string
	cursorPos      int
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
	// Adds 1 to wg counter to prevent console (and printer) from closing befor printer is done
	// Later decremented in Close()
	console.printHandlerWG.Add(1)
	console.CommandManager = manager.NewCommandManager(console)
	console.CommandManager.Init()
	return console
}

func (console *Console) Run() {
	//To remove
	console.consoleWaitGroup.Add(3)
	go console.handleInput()
	go console.printHandler()

	//To remove
	console.consoleWaitGroup.Wait()
}

func (console *Console) Close() {
	console.mutex.Lock()
	defer console.mutex.Unlock()

	if !console.isRunning {
		return
	}

	// Explanation is in Init()
	console.printHandlerWG.Done()
	console.isWaitingForInput = false

	go func() {
		console.printHandlerWG.Wait()
		console.isRunning = false
		close(console.printChan)
		keyboard.Close()
		console.mainWaitGroup.Done()
	}()
}

func (console *Console) Pause() {
	console.mutex.Lock()
	defer console.mutex.Unlock()

	if console.isPaused {
		return
	}

	if console.isWaitingForInput {
		fmt.Print("\r\033[K")
	}

	console.isPaused = true
	console.isWaitingForInput = false

	keyboard.Close()
}

func (console *Console) Resume() {
	console.mutex.Lock()
	defer console.mutex.Unlock()

	if !console.isPaused {
		return
	}

	// Reopen keyboard
	if err := keyboard.Open(); err != nil {
		handler.HandleError(console, "Failed to resume console", err, console.service, true)
		return
	}

	for _, msg := range console.pendingMessages {
		console.printHandlerWG.Done()
		if console.isClosing || !console.isRunning {
			return
		}

		if console.isWaitingForInput {
			fmt.Print("\r\033[K")
		}

		fmt.Println(msg)

		if console.isWaitingForInput {
			console.printPrompt()
		}
	}

	console.isPaused = false
	console.isWaitingForInput = true

	console.printPrompt()
}

func (console *Console) CloseGracefully(message string) {
	console.SetClosingStatus()
	console.Print(message, console.Service())
	console.Close()
}

func (console *Console) Service() service.Service {
	return console.service
}

func (console *Console) handleInput() {
	defer console.consoleWaitGroup.Done()
	if err := keyboard.Open(); err != nil {
		handler.HandleError(console, "App is shutting down...", err, console.service, true)
	}

	console.isWaitingForInput = true
	console.printPrompt()

	for console.isRunning {

		console.mutex.RLock()
		if console.isPaused {
			console.mutex.RUnlock()
			continue
		}
		console.mutex.RUnlock()

		char, key, err := keyboard.GetKey()

		if console.isClosing || !console.isRunning {
			return
		}

		console.mutex.RLock()
		isPaused := console.isPaused
		console.mutex.RUnlock()

		if err != nil {
			if isPaused {
				continue
			}

			console.PrintColored(err.Error(), console.service, model.Red)
			continue
		}

		switch key {
		// If arrows do not work change terminal to xterm-256color 'screen -T xterm-256color -S netmanager'
		// There is a problem with it on screen
		case keyboard.KeyArrowLeft:
			console.moveCursor(-1)
		case keyboard.KeyArrowRight:
			console.moveCursor(1)
		case keyboard.KeyEnter:
			console.executeCommand()
		case keyboard.KeyBackspace, keyboard.KeyBackspace2:
			console.handleBackspace()
		case keyboard.KeySpace:
			console.handleChar(' ')
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
		go console.CommandManager.ExecuteCommand(command)
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
	fmt.Print("\r\033[K")
	fmt.Printf("\r> %s", console.inputBuffer)
	fmt.Printf("\033[%dG", console.cursorPos+3)
}

func (console *Console) moveCursor(delta int) {
	newPos := console.cursorPos + delta
	if newPos < 0 || newPos > len(console.inputBuffer) {
		return
	}
	console.cursorPos = newPos
	fmt.Printf("\033[%dG", console.cursorPos+3)
}

func (console *Console) printHandler() {
	defer console.consoleWaitGroup.Done()

	for msg := range console.printChan {
		console.mutex.RLock()
		isPaused := console.isPaused
		console.mutex.RUnlock()

		if isPaused {
			console.mutex.Lock()
			console.pendingMessages = append(console.pendingMessages, msg)
			console.mutex.Unlock()
			continue
		}

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
		console.printHandlerWG.Done()
		console.mutex.Unlock()
	}
}

func (console *Console) printPrompt() {
	fmt.Print("> ")
}

func (console *Console) Print(message string, service service.Service) {
	console.printHandlerWG.Add(1)
	console.printChan <- fmt.Sprintf("[%s | %s]: %s", time.Now().Format("15:04:05"), service.Name, message)
}

func (console *Console) PrintColored(message string, service service.Service, color int) {
	console.printHandlerWG.Add(1)
	coloredMessage := fmt.Sprintf("\033[38;5;%dm[%s | %s]: %s\033[0m", color, time.Now().Format("15:04:05"), service.Name, message)
	console.printChan <- coloredMessage
}

func (console *Console) SetClosingStatus() {
	console.isClosing = true
}
