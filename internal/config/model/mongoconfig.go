package model

import (
	"NetManager/internal/cli/handler"
	"NetManager/pkg/interfaces"
	"encoding/json"
	"os"
	"path/filepath"
)

type MongoConfig struct {
	ServiceUsername string `json:"service-username"`
	ServicePassword string `json:"service-password"`
	InternalURI     string `json:"internal-uri"`
	ExternalURI     string `json:"external-uri"`
}

func NewDefaultMongoConfig() *MongoConfig {
	return &MongoConfig{
		ServiceUsername: "netmanager",
		ServicePassword: "YM6EMqGQ5NyB9qp3ovkn",
		InternalURI:     "",
		ExternalURI:     "",
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

func (config *MongoConfig) GetServiceUsername() string {
	return config.ServiceUsername
}

func (config *MongoConfig) GetServicePassword() string {
	return config.ServicePassword
}

func (config *MongoConfig) GetExternalURI() string {
	return config.ExternalURI
}

func (config *MongoConfig) GetInternalURI() string {
	return config.InternalURI
}
