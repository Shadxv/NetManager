package model

import (
	"NetManager/internal/cli/handler"
	"NetManager/pkg/interfaces"
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

func CreateDefaultHarborConfigFile(printer interfaces.Printer) *HarborConfig {
	config := NewDefaultHarborConfig()
	jsonData, _ := json.MarshalIndent(config, "", " ")

	filePath := filepath.Join("config", "harbor.json")
	err := os.WriteFile(filePath, jsonData, 0644)
	handler.HandleError(printer, "Error occured during loading Harbor config.", err, printer.Service(), true)

	return config
}

func (config *HarborConfig) GetDomain() string {
	return config.Domain
}

func (config *HarborConfig) GetProjectName() string {
	return config.ProjectName
}

func (config *HarborConfig) GetUsername() string {
	return config.Username
}

func (config *HarborConfig) GetPassword() string {
	return config.Password
}
