package model

import (
	"NetManager/external/cli"
	"NetManager/internal/cli/handler"
	"encoding/json"
	"os"
	"path/filepath"
)

type HarborConfig struct {
	Domain      string `json:"domain"`
	ProjectName string `json:"project-name"`
	Username    string `json:"username"`
	Password    string `json:"user-password"`
}

func NewDefaultHarborConfig() *HarborConfig {
	return &HarborConfig{
		Domain:      "registry.netmanager.net",
		ProjectName: "netmanager-project",
		Username:    "netmanager",
		Password:    "password",
	}
}

func CreateDefaultHarborConfigFile(printer cli.Printer) *HarborConfig {
	config := NewDefaultHarborConfig()
	jsonData, _ := json.MarshalIndent(config, "", " ")

	filePath := filepath.Join("config", "harbor.json")
	err := os.WriteFile(filePath, jsonData, 0644)
	handler.HandleError(printer, "Error occured during loading config.", err, printer.Service(), true)

	return config
}
