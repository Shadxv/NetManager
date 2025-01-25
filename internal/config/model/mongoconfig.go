package model

import (
	"NetManager/internal/cli/handler"
	"NetManager/pkg/interfaces"
	"encoding/json"
	"os"
	"path/filepath"
)

type MongoConfig struct {
	Port          int    `json:"port"`
	ExternalPort  int    `json:"external-port"`
	Username      string `json:"username"`
	Password      string `json:"password"`
	AuthRequired  bool   `json:"auth-required"`
	Authorization bool   `json:"authorization"`
}

func NewDefaultMongoConfig() *MongoConfig {
	return &MongoConfig{
		Port:          27017,
		ExternalPort:  30004,
		Username:      "netmanager",
		Password:      "YM6EMqGQ5NyB9qp3ovkn",
		AuthRequired:  true,
		Authorization: true,
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

func (config *MongoConfig) GetUsername() string {
	return config.Username
}

func (config *MongoConfig) GetPassword() string {
	return config.Password
}

func (config *MongoConfig) IsAuthRequired() bool {
	return config.AuthRequired
}

func (config *MongoConfig) NeedsAuthorization() bool {
	return config.Authorization
}
