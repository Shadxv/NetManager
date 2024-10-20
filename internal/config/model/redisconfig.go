package model

import (
	"NetManager/internal/cli/handler"
	"NetManager/internal/cli/model"
	"encoding/json"
	"os"
	"path/filepath"
)

type RedisConfig struct {
	DockerImage     string `json:"docker-image"`
	Version         string `json:"version"`
	Port            int    `json:"port"`
	Username        string `json:"username"`
	Password        string `json:"password"`
	MaxMemory       string `json:"max-memory"`
	MaxMemoryPolicy string `json:"max-memory-policy"`
	Timeout         int    `json:"timeout"`
}

func NewDefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		DockerImage:     "redis",
		Version:         "alpine",
		Port:            6379,
		Username:        "default",
		Password:        "",
		MaxMemory:       "256mb",
		MaxMemoryPolicy: "allkeys-lru",
		Timeout:         300,
	}
}

func CreateDefaultRedisConfig(printer model.Printer) RedisConfig {
	config := NewDefaultRedisConfig()
	jsonData, _ := json.MarshalIndent(config, "", "  ")

	filePath := filepath.Join("config", "redis.json")

	err := os.WriteFile(filePath, jsonData, 0644)
	handler.HandleError(printer, "Error occured during loading redis config.", err, printer.Service(), true)

	return *config
}
