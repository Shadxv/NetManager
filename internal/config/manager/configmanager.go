package manager

import (
	"NetManager/internal/cli/handler"
	configModel "NetManager/internal/config/model"
	"NetManager/internal/module"
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"NetManager/pkg/util"
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

type ConfigManager struct {
	printer      interfaces.Printer `gob:"-"`
	MainConfig   *configModel.MainConfig
	RedisConfig  *configModel.RedisConfig
	HarborConfig *configModel.HarborConfig
	MongoConfig  *configModel.MongoConfig
	status       string `gob:"-"`
}

func (configManager *ConfigManager) Init(moduleManager *module.Manager) {

}

func (configManager *ConfigManager) Disable(shutdown bool) {
	if !shutdown {
		if configManager.printer != nil {
			configManager.printer.PrintColored("This module cannot be disabled!", configManager.printer.Service(), types.Red)
		}
		return
	}
	configManager.status = types.Stopping
	// save data (services)
	configManager.status = types.Disabled
}

func (configManager *ConfigManager) Reload() {
	startTime := time.Now()
	if configManager.printer != nil {
		configManager.printer.Print("Reloading config module...", configManager.printer.Service())
	}
	configManager.LoadData()
	if configManager.printer != nil {
		millisecondsElapsed := time.Since(startTime).Milliseconds()
		configManager.printer.Print("Reloaded config module in "+strconv.FormatInt(millisecondsElapsed, 10)+"ms.", configManager.printer.Service())
	}
}

func (configManager *ConfigManager) SaveData() error {
	return util.SaveData(configManager.Type(), configManager)
}

func (configManager *ConfigManager) LoadData() {
	configManager.loadFolderStructure()
	configManager.loadMainConfigFile()
	configManager.loadRedisConfigFile()
	configManager.loadHarborConfigFile()
	configManager.loadMongoConfigFile()
}

func (configManager *ConfigManager) SetStatus(newStatus string) {
	configManager.status = newStatus
}

func (configManager *ConfigManager) Status() string {
	return configManager.status
}

func (configManager *ConfigManager) Type() string {
	return types.Config
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{}
}

func (configManager *ConfigManager) GetMainConfig() interfaces.MainConfig {
	return configManager.MainConfig
}

func (configManager *ConfigManager) GetRedisConfig() interfaces.RedisConfig {
	return configManager.RedisConfig
}

func (configManager *ConfigManager) GetHarborConfig() interfaces.HarborConfig {
	return configManager.HarborConfig
}

func (configManager *ConfigManager) GetMongoConfig() interfaces.MongoConfig {
	return configManager.MongoConfig
}

func (configManager *ConfigManager) loadFolderStructure() {
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
		configManager.MainConfig = configModel.CreateDefaultMainConfig(configManager.printer)
		return
	}

	fileData, err := os.ReadFile(filePath)

	if handler.HandleError(configManager.printer, "Error occured during loading config.", err, configManager.printer.Service(), true) {
		return
	}

	err = json.Unmarshal(fileData, &configManager.MainConfig)
	if handler.HandleError(configManager.printer, "Error occured during loading config.", err, configManager.printer.Service(), true) {
		return
	}
}

func (configManager *ConfigManager) loadRedisConfigFile() {
	filePath := filepath.Join("config", "redis.json")
	configManager.printer.Print("Loading redis config file...", configManager.printer.Service())

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		configManager.RedisConfig = configModel.CreateDefaultRedisConfig(configManager.printer)
		return
	}

	fileData, err := os.ReadFile(filePath)

	if handler.HandleError(configManager.printer, "Error occured during loading Redis config.", err, configManager.printer.Service(), true) {
		return
	}

	err = json.Unmarshal(fileData, &configManager.RedisConfig)
	if handler.HandleError(configManager.printer, "Error occured during loading Redis config.", err, configManager.printer.Service(), true) {
		return
	}
}

func (configManager *ConfigManager) loadHarborConfigFile() {
	filePath := filepath.Join("config", "harbor.json")
	configManager.printer.Print("Loading Harbor config file...", configManager.printer.Service())

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		configManager.HarborConfig = configModel.CreateDefaultHarborConfigFile(configManager.printer)
		return
	}

	fileData, err := os.ReadFile(filePath)

	if handler.HandleError(configManager.printer, "Error occured during loading Harbor config.", err, configManager.printer.Service(), true) {
		return
	}

	err = json.Unmarshal(fileData, &configManager.HarborConfig)
	if handler.HandleError(configManager.printer, "Error occured during loading Harbor config.", err, configManager.printer.Service(), true) {
		return
	}
}

func (configManager *ConfigManager) loadMongoConfigFile() {
	filePath := filepath.Join("config", "mongodb.json")
	configManager.printer.Print("Loading MongoDB config file...", configManager.printer.Service())

	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		configManager.MongoConfig = configModel.CreateDefaultMongoConfig(configManager.printer)
		return
	}

	fileData, err := os.ReadFile(filePath)

	if handler.HandleError(configManager.printer, "Error occured during loading MongoDB config.", err, configManager.printer.Service(), true) {
		return
	}

	err = json.Unmarshal(fileData, &configManager.MongoConfig)
	if handler.HandleError(configManager.printer, "Error occured during loading MongoDB config.", err, configManager.printer.Service(), true) {
		return
	}
}
