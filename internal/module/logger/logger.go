package logger

import (
	"bufio"
	"os"
	"strings"
)

const LogPath = "data/log.txt"

type Logger struct {
	saveChannel chan string
	file        *os.File
	fileExists  bool
	scanner     *bufio.Scanner
	writer      *bufio.Writer
}

func (logger *Logger) Init() {
	if instance != nil {
		return
	}

	instance = logger

	file, err := os.OpenFile(LogPath, os.O_RDONLY, 0644)
	logger.fileExists = err == nil
	if logger.fileExists {
		logger.file = file
	}
}

func (logger *Logger) IsPresent() bool {
	return logger.fileExists
}

func (logger *Logger) Finalize() error {
	logger.scanner = nil
	logger.file.Close()

	if logger.fileExists {
		err := os.Remove(LogPath)
		if err != nil {
			return err
		}
	}

	file, err := os.OpenFile(LogPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	logger.file = file
	logger.writer = bufio.NewWriter(logger.file)
	logger.saveChannel = make(chan string)

	go func() {
		for log := range logger.saveChannel {
			logger.writer.WriteString(log + "\n")
			logger.writer.Flush()
		}
	}()

	return nil
}

func (logger *Logger) ReadFile() {
	if logger.file == nil {
		return
	}
	logger.scanner = bufio.NewScanner(logger.file)
}

func (logger *Logger) ReadNext() (bool, string, []string) {
	if logger.file == nil || logger.scanner == nil || logger.scanner.Scan() {
		return false, "", nil
	}

	line := logger.scanner.Text()
	parts := strings.SplitN(line, " ", 2)
	if len(parts) != 2 {
		return true, "corrupted", nil
	}
	module, log := parts[0], strings.Split(parts[1], ";")
	return true, module, log
}

func (logger *Logger) SaveLog(module string, data []string) {
	logLine := module + " " + strings.Join(data, ";")
	logger.saveChannel <- logLine
}

func (logger *Logger) Close() {
	instance = nil
	close(logger.saveChannel)
	if logger.file == nil {
		return
	}
	logger.file.Close()
}

var instance *Logger = nil

func GetInstance() *Logger {
	if instance == nil {
		return &Logger{}
	}
	return instance
}
