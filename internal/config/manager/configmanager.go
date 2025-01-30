package manager

import (
	"NetManager/internal/cli/handler"
	configModel "NetManager/internal/config/model"
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"encoding/json"
	"os"
	"path"
	"path/filepath"
)

type ConfigManager struct {
	printer      interfaces.Printer
	mainConfig   *configModel.MainConfig
	redisConfig  *configModel.RedisConfig
	harborConfig *configModel.HarborConfig
	mongoConfig  *configModel.MongoConfig
}

func NewConfigManager(printer interfaces.Printer) *ConfigManager {
	return &ConfigManager{
		printer: printer,
	}
}

func (configManager *ConfigManager) GetMainConfig() interfaces.MainConfig {
	return configManager.mainConfig
}

func (configManager *ConfigManager) GetRedisConfig() interfaces.RedisConfig {
	return configManager.redisConfig
}

func (configManager *ConfigManager) GetHarborConfig() interfaces.HarborConfig {
	return configManager.harborConfig
}

func (configManager *ConfigManager) GetMongoConfig() interfaces.MongoConfig {
	return configManager.mongoConfig
}

func (configManager *ConfigManager) Init() {
	configManager.loadFolderStructure()
	configManager.loadMainConfigFile()
	configManager.loadRedisConfigFile()
	configManager.loadHarborConfigFile()
	configManager.loadMongoConfigFile()
}

func (configManager *ConfigManager) loadFolderStructure() {
	configManager.loadFolder("config", "")
	configManager.loadFolder("data", "")
	configManager.loadFolder("services", "")
	configManager.loadFolder("templates", "services")
	configManager.loadFolder("instances", "services")
	configManager.loadFolder("config", "services")
	configManager.loadFolder("paper-default", path.Join("services", "templates"))
	configManager.loadFolder("velocity-default", path.Join("services", "templates"))
}

func (configManager *ConfigManager) loadFolder(name string, path string) {
	fullPath := filepath.Join(path, name)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		configManager.printer.Print("'"+name+"' folder not found. Creating new one...", configManager.printer.Service())
		err := os.Mkdir(fullPath, 0755)
		if err != nil {
			configManager.printer.PrintColored("Could not create new '"+name+"' folder...", configManager.printer.Service(), types.Red)
			configManager.printer.PrintColored(err.Error(), configManager.printer.Service(), types.Red)
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

	if handler.HandleError(configManager.printer, "Error occured during loading Redis config.", err, configManager.printer.Service(), true) {
		return
	}

	err = json.Unmarshal(fileData, &configManager.redisConfig)
	if handler.HandleError(configManager.printer, "Error occured during loading Redis config.", err, configManager.printer.Service(), true) {
		return
	}
}

func (configManager *ConfigManager) loadHarborConfigFile() {
	filePath := filepath.Join("config", "harbor.json")
	configManager.printer.Print("Loading Harbor config file...", configManager.printer.Service())

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		configManager.harborConfig = configModel.CreateDefaultHarborConfigFile(configManager.printer)
		return
	}

	fileData, err := os.ReadFile(filePath)

	if handler.HandleError(configManager.printer, "Error occured during loading Harbor config.", err, configManager.printer.Service(), true) {
		return
	}

	err = json.Unmarshal(fileData, &configManager.harborConfig)
	if handler.HandleError(configManager.printer, "Error occured during loading Harbor config.", err, configManager.printer.Service(), true) {
		return
	}
}

func (configManager *ConfigManager) loadMongoConfigFile() {
	filePath := filepath.Join("config", "mongodb.json")
	configManager.printer.Print("Loading MongoDB config file...", configManager.printer.Service())

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		configManager.mongoConfig = configModel.CreateDefaultMongoConfig(configManager.printer)
		return
	}

	fileData, err := os.ReadFile(filePath)

	if handler.HandleError(configManager.printer, "Error occured during loading MongoDB config.", err, configManager.printer.Service(), true) {
		return
	}

	err = json.Unmarshal(fileData, &configManager.mongoConfig)
	if handler.HandleError(configManager.printer, "Error occured during loading MongoDB config.", err, configManager.printer.Service(), true) {
		return
	}
}
