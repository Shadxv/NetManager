package model

import (
	"NetManager/external/cli"
	"NetManager/internal/cli/handler"
	"encoding/json"
	"os"
	"path/filepath"
)

type MainConfig struct {
	Name string `json:"service-name"`
}

func NewDefaultMainConfig() *MainConfig {
	return &MainConfig{
		Name: "netmanager-service",
	}
}

func CreateDefaultMainConfig(printer cli.Printer) MainConfig {
	config := NewDefaultMainConfig()
	jsonData, _ := json.MarshalIndent(config, "", "  ")

	filePath := filepath.Join("config", "config.json")

	err := os.WriteFile(filePath, jsonData, 0644)
	handler.HandleError(printer, "Error occured during loading config.", err, printer.Service(), true)

	return *config
}
