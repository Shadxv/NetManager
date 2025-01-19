package model

import (
	"NetManager/pkg/interfaces"
)

type VelocityData struct {
	version        string
	build          int
	port           int
	replicasAmount int
}

func NewVelocityData(version string, build int, port int, replicasAmount int) *VelocityData {
	return &VelocityData{
		version:        version,
		build:          build,
		port:           port,
		replicasAmount: replicasAmount,
	}
}

func (data *VelocityData) Version() string {
	return data.version
}

func (data *VelocityData) BuildNumber() int {
	return data.build
}

func (data *VelocityData) Port() int {
	return data.port
}

func (data *VelocityData) ReplicasAmount() int {
	return data.replicasAmount
}

func (data *VelocityData) Build(serviceModel interfaces.ServiceModel, printer interfaces.Printer) {

}

func (data *VelocityData) Update(serviceModel interfaces.ServiceModel, printer interfaces.Printer) {

}

func (data *VelocityData) Stop(serviceModel interfaces.ServiceModel, printer interfaces.Printer) {

}

func (data *VelocityData) Start(serviceModel interfaces.ServiceModel, printer interfaces.Printer) {

}

func (data *VelocityData) Deploy(serviceModel interfaces.ServiceModel, printer interfaces.Printer) {

}
