package model

import (
	"NetManager/external/cli"
	"NetManager/internal/cli/handler"
	"encoding/json"
	"os"
	"path/filepath"
)

type HarborConfig struct {
	httpPort     int    `json:"http-port"`
	projectName  string `json:"project-name"`
	username     string `json:"username"`
	userMail     string `json:"user-mail"`
	userPassword string `json:"user-password"`
	userRole     string `json:"user-role"`
	disableGuest bool   `json:"disable-guest"`
}

func NewDefaultHarborConfig() *HarborConfig {
	return &HarborConfig{
		httpPort: 8001,

		projectName: "netmanager-project",
		username:    "netmanage",
		userMail:    "netmanager@netmanager.net",
		userPassword: "" +
			"",
		userRole:     "maintainer",
		disableGuest: true,
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

func (config HarborConfig) HttpPort() int {
	return config.httpPort
}

func (config HarborConfig) ProjectName() string {
	return config.projectName
}

func (config HarborConfig) Username() string {
	return config.username
}

func (config HarborConfig) UserMail() string {
	return config.userMail
}

func (config HarborConfig) UserPassword() string {
	return config.userPassword
}

func (config HarborConfig) UserRole() string {
	return config.userRole
}

func (config HarborConfig) DisableGuest() bool {
	return config.disableGuest
}
