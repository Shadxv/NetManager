package model

import (
	"NetManager/internal/cli/handler"
	"NetManager/pkg/interfaces"
	"encoding/json"
	"os"
	"path/filepath"
)

type MongoConfig struct {
	Port            int    `json:"port"`
	ExternalPort    int    `json:"external-port"`
	RootUsername    string `json:"root-username"`
	RootPassword    string `json:"root-password"`
	ServiceUsername string `json:"service-username"`
	ServicePassword string `json:"service-password"`
	AuthRequired    bool   `json:"auth-required"`
}

func NewDefaultMongoConfig() *MongoConfig {
	return &MongoConfig{
		Port:            27017,
		ExternalPort:    30004,
		RootUsername:    "admin",
		RootPassword:    "admin",
		ServiceUsername: "netmanager",
		ServicePassword: "YM6EMqGQ5NyB9qp3ovkn",
		AuthRequired:    true,
	}
}

func CreateDefaultMongoConfig(printer interfaces.Printer) *MongoConfig {
	config := NewDefaultMongoConfig()
	jsonData, _ := json.MarshalIndent(config, "", "  ")

	filePath := filepath.Join("config", "mongodb.json")

	err := os.WriteFile(filePath, jsonData, 0644)
	handler.HandleError(printer, "Error occured during loading MongoDB config.", err, printer.Service(), true)

	return config
}

func (config *MongoConfig) GetPort() int {
	return config.Port
}

func (config *MongoConfig) GetExternalPort() int {
	return config.ExternalPort
}

func (config *MongoConfig) GetRootUsername() string {
	return config.RootUsername
}

func (config *MongoConfig) GetRootPassword() string {
	return config.RootPassword
}

func (config *MongoConfig) GetServiceUsername() string {
	return config.ServiceUsername
}

func (config *MongoConfig) GetServicePassword() string {
	return config.ServicePassword
}

func (config *MongoConfig) IsAuthRequired() bool {
	return config.AuthRequired
}
