package manager

import (
	cliModel "NetManager/internal/cli/model"
	configModel "NetManager/internal/config/model"
	"encoding/json"
	"os"
	"path/filepath"
)

type ConfigManager struct {
	printer    cliModel.Printer
	mainConfig configModel.MainConfig
}

func NewConfigManager(printer cliModel.Printer) *ConfigManager {
	return &ConfigManager{
		printer: printer,
	}
}

func (configManager *ConfigManager) GetMainConfig() *configModel.MainConfig {
	return &configManager.mainConfig
}

func (configManager *ConfigManager) Init() {
	configManager.loadFolderStructure()
	configManager.loadMainConfigFile()
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

	if err != nil {
		configManager.printer.PrintColored("Error occured during loading config.", configManager.printer.Service(), cliModel.Red)
		configManager.printer.PrintColored(err.Error(), configManager.printer.Service(), cliModel.Red)
		configManager.printer.CloseGracefully("App is shutting down...")
		return
	}

	err = json.Unmarshal(fileData, &configManager.mainConfig)
	if err != nil {
		configManager.printer.PrintColored("Error occured during loading config.", configManager.printer.Service(), cliModel.Red)
		configManager.printer.PrintColored(err.Error(), configManager.printer.Service(), cliModel.Red)
		configManager.printer.CloseGracefully("App is shutting down...")
		return
	}
}
