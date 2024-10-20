package manager

import (
	"NetManager/internal/cli/handler"
	cliModel "NetManager/internal/cli/model"
	configModel "NetManager/internal/config/model"
	"encoding/json"
	"os"
	"path/filepath"
)

type ConfigManager struct {
	printer     cliModel.Printer
	mainConfig  configModel.MainConfig
	redisConfig configModel.RedisConfig
}

func NewConfigManager(printer cliModel.Printer) *ConfigManager {
	return &ConfigManager{
		printer: printer,
	}
}

func (configManager *ConfigManager) GetMainConfig() *configModel.MainConfig {
	return &configManager.mainConfig
}

func (configManager *ConfigManager) GetRedisConfig() *configModel.RedisConfig {
	return &configManager.redisConfig
}

func (configManager *ConfigManager) Init() {
	configManager.loadFolderStructure()
	configManager.loadMainConfigFile()
	configManager.loadRedisConfigFile()
}

func (configManager *ConfigManager) loadFolderStructure() {
	configManager.loadFolder("config", "")
	configManager.loadFolder("data", "")
}

func (configManager *ConfigManager) loadFolder(name string, path string) {
	fullPath := filepath.Join(path, name)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		configManager.printer.Print("'"+name+"' folder not found. Creating new one...", configManager.printer.Service())
		err := os.Mkdir(fullPath, 0755)
		if err != nil {
			configManager.printer.PrintColored("Could not create new '"+name+"' folder...", configManager.printer.Service(), cliModel.Red)
			configManager.printer.PrintColored(err.Error(), configManager.printer.Service(), cliModel.Red)
			configManager.printer.CloseGracefully("App is shutting down...")
		}
	}
}

func (configManager *ConfigManager) loadMainConfigFile() {
	filePath := filepath.Join("config", "config.json")
	configManager.printer.Print("Loading config file...", configManager.printer.Service())

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		configManager.mainConfig = configModel.CreateDefaultMainConfig(configManager.printer)
		return
	}

	fileData, err := os.ReadFile(filePath)

	if handler.HandleError(configManager.printer, "Error occured during loading config.", err, configManager.printer.Service(), true) {
		return
	}

	err = json.Unmarshal(fileData, &configManager.mainConfig)
	if handler.HandleError(configManager.printer, "Error occured during loading config.", err, configManager.printer.Service(), true) {
		return
	}
}

func (configManager *ConfigManager) loadRedisConfigFile() {
	filePath := filepath.Join("config", "redis.json")
	configManager.printer.Print("Loading redis config file...", configManager.printer.Service())

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		configManager.redisConfig = configModel.CreateDefaultRedisConfig(configManager.printer)
		return
	}

	fileData, err := os.ReadFile(filePath)

	if handler.HandleError(configManager.printer, "Error occured during loading redis config.", err, configManager.printer.Service(), true) {
		return
	}

	err = json.Unmarshal(fileData, &configManager.redisConfig)
	if handler.HandleError(configManager.printer, "Error occured during loading redis config.", err, configManager.printer.Service(), true) {
		return
	}
}
