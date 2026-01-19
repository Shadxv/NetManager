package model

import (
	"NetManager/internal/cli/handler"
	"NetManager/pkg/interfaces"
	"encoding/json"
	"os"
	"path/filepath"
)

type MainConfig struct {
	Name            string `json:"service-name"`
	ServerGroupName string `json:"server-group-name"`
	JWTSecret       string `json:"jwt-secret"`
}

func NewDefaultMainConfig() *MainConfig {
	return &MainConfig{
		Name:            "netmanager-service",
		ServerGroupName: "dreammc",
		JWTSecret:       "change-me-please-to-something-secure",
	}
}

func CreateDefaultMainConfig(printer interfaces.Printer) *MainConfig {
	config := NewDefaultMainConfig()
	jsonData, _ := json.MarshalIndent(config, "", "  ")

	filePath := filepath.Join("config", "config.json")

	err := os.WriteFile(filePath, jsonData, 0644)
	handler.HandleError(printer, "Error occured during loading config.", err, printer.Service(), true)

	return config
}

func (config *MainConfig) GetName() string {
	return config.Name
}

func (config *MainConfig) GetServerGroupName() string {
	return config.ServerGroupName
}

func (config *MainConfig) GetJWTSecret() string {
	return config.JWTSecret
}
