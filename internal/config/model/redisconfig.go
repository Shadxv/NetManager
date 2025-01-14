package model

import (
	"NetManager/external/cli"
	"NetManager/internal/cli/handler"
	"encoding/json"
	"os"
	"path/filepath"
)

type RedisConfig struct {
	DockerImage     string `json:"docker-image"`
	Version         string `json:"version"`
	Port            int    `json:"port"`
	Password        string `json:"password"`
	MaxMemory       string `json:"max-memory"`
	MaxMemoryPolicy string `json:"max-memory-policy"`
}

func NewDefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		DockerImage:     "redis",
		Version:         "alpine",
		Port:            6379,
		Password:        "",
		MaxMemory:       "256mb",
		MaxMemoryPolicy: "allkeys-lru",
	}
}

func CreateDefaultRedisConfig(printer cli.Printer) *RedisConfig {
	config := NewDefaultRedisConfig()
	jsonData, _ := json.MarshalIndent(config, "", "  ")

	filePath := filepath.Join("config", "redis.json")

	err := os.WriteFile(filePath, jsonData, 0644)
	handler.HandleError(printer, "Error occured during loading redis config.", err, printer.Service(), true)

	return config
}
