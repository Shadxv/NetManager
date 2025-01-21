package model

import (
	"NetManager/internal/cli/handler"
	"NetManager/pkg/interfaces"
	"encoding/json"
	"os"
	"path/filepath"
)

type RedisConfig struct {
	DockerImage     string `json:"docker-image"`
	Version         string `json:"version"`
	Port            int    `json:"port"`
	ExternalPort    int    `json:"external-port"`
	Password        string `json:"password"`
	MaxMemory       string `json:"max-memory"`
	MaxMemoryPolicy string `json:"max-memory-policy"`
}

func NewDefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		DockerImage:     "redis",
		Version:         "alpine",
		Port:            6378,
		ExternalPort:    31999,
		Password:        "",
		MaxMemory:       "256mb",
		MaxMemoryPolicy: "allkeys-lru",
	}
}

func CreateDefaultRedisConfig(printer interfaces.Printer) *RedisConfig {
	config := NewDefaultRedisConfig()
	jsonData, _ := json.MarshalIndent(config, "", "  ")

	filePath := filepath.Join("config", "redis.json")

	err := os.WriteFile(filePath, jsonData, 0644)
	handler.HandleError(printer, "Error occured during loading redis config.", err, printer.Service(), true)

	return config
}

func (config *RedisConfig) GetDockerImage() string {
	return config.DockerImage
}

func (config *RedisConfig) GetVersion() string {
	return config.Version
}

func (config *RedisConfig) GetPort() int {
	return config.Port
}

func (config *RedisConfig) GetExternalPort() int {
	return config.ExternalPort
}

func (config *RedisConfig) GetPassword() string {
	return config.Password
}

func (config *RedisConfig) GetMaxMemory() string {
	return config.MaxMemory
}

func (config *RedisConfig) GetMaxMemoryPolicy() string {
	return config.MaxMemoryPolicy
}
