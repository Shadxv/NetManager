package model

import (
	"NetManager/pkg/interfaces"
)

type PaperData struct {
	version     string
	build       int
	minReplicas int
}

func NewPaperData(version string, build int, minReplicas int) *PaperData {
	return &PaperData{
		version:     version,
		build:       build,
		minReplicas: minReplicas,
	}
}

func (data *PaperData) Version() string {
	return data.version
}

func (data *PaperData) BuildNumber() int {
	return data.build
}

func (data *PaperData) MinReplicas() int {
	return data.minReplicas
}

func (data *PaperData) Build(serviceModel interfaces.ServiceModel, printer interfaces.Printer) {
	//imageManager.FullDeployImage(serviceModel, serviceManager)
}

func (data *PaperData) Update(serviceModel interfaces.ServiceModel, printer interfaces.Printer) {

}

func (data *PaperData) Stop(serviceModel interfaces.ServiceModel, printer interfaces.Printer) {

}

func (data *PaperData) Start(serviceModel interfaces.ServiceModel, printer interfaces.Printer) {

}

func (data *PaperData) Deploy(serviceModel interfaces.ServiceModel, printer interfaces.Printer) {

}
