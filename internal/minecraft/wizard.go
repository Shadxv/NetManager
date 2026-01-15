package minecraft

import (
	minecraftModel "NetManager/internal/minecraft/model"
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	Finished = iota
	Canceled
	Error
)

type ServiceWizard struct {
	printer         interfaces.Printer
	serviceManager  interfaces.ServiceManager
	groupName       string
	serviceName     string
	fetchedVersions []string
	fetchedBuilds   []int
	serviceType     string
}

func NetServiceWizard(printer interfaces.Printer, serviceManager interfaces.ServiceManager, serviceName string, groupName string) *ServiceWizard {
	return &ServiceWizard{
		printer:        printer,
		serviceManager: serviceManager,
		groupName:      groupName,
		serviceName:    serviceName,
	}
}

func (wizard *ServiceWizard) Run() {
	if wizard.serviceManager.Exists(wizard.serviceName) {
		wizard.printer.PrintColored("Service with that name already exists.", wizard.printer.Service(), types.Red)
		return
	}
	wizard.init()
	wizard.runConfiguration()
	wizard.close()
}

func (wizard *ServiceWizard) init() {
	wizard.printer.Pause()
	fmt.Println("Starting service wizard for '" + wizard.serviceName + "'...")
}

func (wizard *ServiceWizard) runConfiguration() {
	if !wizard.readServiceType() {
		fmt.Println("Incorrect type. Closing wizard")
		return
	}

	var result int
	var serviceData interface{}

	switch wizard.serviceType {
	case types.Paper:
		result, serviceData = wizard.paperConfiguration()
	case types.Velocity:
		result, serviceData = wizard.velocityConfiguration()
	}

	switch result {
	case Error:
		fmt.Println("Unexpected error occured. Closing wizard")
	case Canceled:
		fmt.Println("Configuration canceled. Closing wizard")
	case Finished:
		wizard.finalize(serviceData)
	}

}

func (wizard *ServiceWizard) finalize(serviceData interface{}) {
	wizard.serviceManager.AddNewService(wizard.serviceName, wizard.serviceType, serviceData)
}

func (wizard *ServiceWizard) close() {
	wizard.printer.Resume()
}

func (wizard *ServiceWizard) readServiceType() bool {
	fmt.Print("Service type (PAPER/velocity): ")

	input := strings.ToLower(readInput())
	switch input {
	case "", "p", "paper":
		wizard.serviceType = types.Paper
	case "v", "velocity":
		wizard.serviceType = types.Velocity
	default:
		return false
	}

	return true
}

func (wizard *ServiceWizard) paperConfiguration() (int, *minecraftModel.PaperData) {
	fmt.Println("Fetching available versions... Please wait")
	if !fetchVersions("https://api.papermc.io/v2/projects/paper/", &wizard.fetchedVersions) {
		return Error, nil
	}

	version := wizard.readVersion("Specify papermc version (default: LATEST): ", "Invalid version...", false)
	if version == "EXIT" {
		return Canceled, nil
	}

	fmt.Println("Fetching available builds... Please wait")
	if !fetchBuilds("https://api.papermc.io/v2/projects/paper/versions/"+version+"/builds/", &wizard.fetchedBuilds) {
		return Error, nil
	}
	build, res := wizard.readBuild("Specify papermc build number (default: LATEST): ", "Invalid build number...", false)
	if res != Finished {
		return res, nil
	}

	minReplicas, res := wizard.readIntValue("Minimum number of service replicas (default: 1): ", "Invalid number of replicas...", 1, 1, math.MaxInt32, false)
	if res != Finished {
		return res, nil
	}

	mongoURI, redisURI := wizard.getURIs()

	return Finished, minecraftModel.NewPaperData(wizard.groupName, version, build, minReplicas, mongoURI, redisURI)
}

func (wizard *ServiceWizard) velocityConfiguration() (int, *minecraftModel.VelocityData) {
	fmt.Println("Fetching available versions... Please wait")
	if !fetchVersions("https://api.papermc.io/v2/projects/velocity/", &wizard.fetchedVersions) {
		return Error, nil
	}

	version := wizard.readVersion("Specify velocity version (default: LATEST): ", "Invalid version...", false)
	if version == "EXIT" {
		return Canceled, nil
	}

	fmt.Println("Fetching available builds... Please wait")
	if !fetchBuilds("https://api.papermc.io/v2/projects/velocity/versions/"+version+"/builds/", &wizard.fetchedBuilds) {
		return Error, nil
	}
	build, res := wizard.readBuild("Specify velocity build number (default: LATEST): ", "Invalid build number...", false)
	if res != Finished {
		return res, nil
	}

	port, res := wizard.readIntValue("Specify proxy port (default: 25565): ", "Invalid port...", 25565, 1, 65535, false)
	if res != Finished {
		return res, nil
	}

	replicasAmount, res := wizard.readIntValue("Replicas amount (default: 1): ", "Invalid number of replicas...", 1, 1, math.MaxInt32, false)
	if res != Finished {
		return res, nil
	}

	mongoURI, redisURI := wizard.getURIs()

	return Finished, minecraftModel.NewVelocityData(wizard.groupName, version, build, port, replicasAmount, mongoURI, redisURI)
}

func (wizard *ServiceWizard) getURIs() (string, string) {
	mongoService := wizard.serviceManager.GetService("mongodb")
	if mongoService == nil {
		return "", ""
	}

	mongoData := interfaces.GetMongoData(mongoService)
	if mongoData == nil {
		return "", ""
	}

	redisService := wizard.serviceManager.GetService("redis")
	if redisService == nil {
		return "", ""
	}

	redisData := interfaces.GetRedisData(redisService)
	if redisData == nil {
		return "", ""
	}

	return mongoData.InternalURI(), fmt.Sprintf("redis://:%s@%s:%d/", redisData.Password(), redisData.InternalRedisIp(), redisData.Port())
}

func (wizard *ServiceWizard) readVersion(msg string, errMsg string, invalid bool) string {
	if invalid {
		fmt.Println(errMsg)
	}
	fmt.Print(msg)
	input := strings.ToUpper(readInput())
	switch input {
	case "EXIT", "STOP", "CANCEL":
		return "EXIT"
	case "LATEST", "":
		return wizard.fetchedVersions[len(wizard.fetchedVersions)-1]
	default:
		if slices.Contains(wizard.fetchedVersions, input) {
			return input
		}
		return wizard.readVersion(msg, errMsg, true)
	}
}

func (wizard *ServiceWizard) readBuild(msg string, errMsg string, invalid bool) (int, int) {
	if invalid {
		fmt.Println(errMsg)
	}
	fmt.Print(msg)
	input := strings.ToUpper(readInput())
	switch input {
	case "EXIT", "STOP", "CANCEL":
		return 0, Canceled
	case "LATEST", "":
		return wizard.fetchedBuilds[len(wizard.fetchedBuilds)-1], Finished
	default:
		build, err := strconv.Atoi(input)
		if err == nil && slices.Contains(wizard.fetchedBuilds, build) {
			return build, Finished
		}
		return wizard.readBuild(msg, errMsg, true)
	}
}

func (wizard *ServiceWizard) readIntValue(msg string, errMsg string, defaultValue int, minValue int, maxValue int, invalid bool) (int, int) {
	if invalid {
		fmt.Println(errMsg)
	}
	fmt.Print(msg)
	input := strings.ToLower(readInput())
	switch input {
	case "exit", "stop", "cancel":
		return 0, Canceled
	case "":
		return defaultValue, Finished
	default:
		value, err := strconv.Atoi(input)
		if err == nil && minValue <= value && value <= maxValue {
			return value, Finished
		}
		return wizard.readIntValue(msg, errMsg, defaultValue, minValue, maxValue, true)
	}
}

func fetchData(url string) map[string]interface{} {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil
	}

	return result
}

func fetchVersions(url string, data *[]string) bool {
	result := fetchData(url)

	if result == nil {
		return false
	}

	if versions, ok := result["versions"].([]interface{}); ok {
		*data = make([]string, len(versions))
		for i, v := range versions {
			(*data)[i] = v.(string)
		}
		return true
	}

	return false
}

func fetchBuilds(url string, data *[]int) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false
	}

	if builds, ok := result["builds"].([]interface{}); ok {
		*data = make([]int, len(builds))
		for i, build := range builds {
			if buildObj, ok := build.(map[string]interface{}); ok {
				if buildNum, ok := buildObj["build"].(float64); ok {
					(*data)[i] = int(buildNum)
				}
			}
		}
		return true
	}

	return false
}

func readInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}
