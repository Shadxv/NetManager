package model

import (
	"NetManager/internal/cli/model"
	"encoding/json"
	"os"
	"path/filepath"
)

type MainConfig struct {
	Name string `json:"service-name"`
}

func NewDefaultMainConfig() *MainConfig {
	return &MainConfig{
		Name: "NetManager-Server",
	}
}

func CreateDefaultMainConfig(printer model.Printer) MainConfig {
	config := NewDefaultMainConfig()
	jsonData, _ := json.MarshalIndent(config, "", "  ")

	filePath := filepath.Join("config", "config.json")

	err := os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		printer.PrintColored("Error occured during loading config.", printer.Service(), model.Red)
		printer.PrintColored(err.Error(), printer.Service(), model.Red)
		printer.CloseGracefully("App is shutting down...")
	}

	return *config
}
