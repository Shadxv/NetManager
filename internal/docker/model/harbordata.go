package model

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
)

type HarborData struct {
	domain      string
	projectName string
	username    string
	password    string
}

func NewHarborData(domain string, projectName string, username string, password string) *HarborData {
	return &HarborData{
		domain:      domain,
		projectName: projectName,
		username:    username,
		password:    password,
	}
}

func (data *HarborData) Domain() string {
	return data.domain
}

func (data *HarborData) ProjectName() string {
	return data.projectName
}

func (data *HarborData) Username() string {
	return data.username
}

func (data *HarborData) Password() string {
	return data.password
}

func (data *HarborData) Build(serviceModel interfaces.ServiceModel, printer interfaces.Printer) {
	printer.PrintColored("Harbor service cannot be built.", printer.Service(), types.Yellow)
}

func (data *HarborData) Update(serviceModel interfaces.ServiceModel, printer interfaces.Printer) {
	printer.PrintColored("Harbor service cannot be updated.", printer.Service(), types.Yellow)
}

func (data *HarborData) Stop(serviceModel interfaces.ServiceModel, printer interfaces.Printer) {
	printer.PrintColored("Harbor service cannot be stopped.", printer.Service(), types.Yellow)
}

func (data *HarborData) Start(serviceModel interfaces.ServiceModel, printer interfaces.Printer) {
	printer.PrintColored("Harbor service cannot be started.", printer.Service(), types.Yellow)
}

func (data *HarborData) Deploy(serviceModel interfaces.ServiceModel, printer interfaces.Printer) {
	printer.PrintColored("Harbor service cannot be deployed.", printer.Service(), types.Yellow)
}
