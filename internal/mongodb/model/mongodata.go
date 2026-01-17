package model

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
)

type MongoData struct {
	ServiceUsernameField string
	ServicePasswordField string
	InternalURIField     string
	ExternalURIField     string
}

func NewMongoData(serviceUsername string, servicePassword string, internalURI string, externalURI string) *MongoData {
	return &MongoData{
		ServiceUsernameField: serviceUsername,
		ServicePasswordField: servicePassword,
		InternalURIField:     internalURI,
		ExternalURIField:     externalURI,
	}
}

func (data *MongoData) ServiceUsername() string {
	return data.ServiceUsernameField
}

func (data *MongoData) ServicePassword() string {
	return data.ServicePasswordField
}

func (data *MongoData) InternalURI() string {
	return data.InternalURIField
}

func (data *MongoData) ExternalURI() string {
	return data.ExternalURIField
}

func (data *MongoData) Build(serviceModel interfaces.ServiceModel, printer interfaces.Printer, imageManager interfaces.ImageManager, serviceManager interfaces.ServiceManager) {
	printer.PrintColored("MongoDB service cannot be built.", printer.Service(), types.Yellow)
}

func (data *MongoData) Update(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	printer.PrintColored("MongoDB service cannot be updated.", printer.Service(), types.Yellow)
}

func (data *MongoData) Stop(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	printer.PrintColored("MongoDB service cannot be stopped.", printer.Service(), types.Yellow)
}

func (data *MongoData) Start(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	printer.PrintColored("MongoDB service cannot be started.", printer.Service(), types.Yellow)
}

func (data *MongoData) Deploy(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	printer.PrintColored("MongoDB service cannot be deployed.", printer.Service(), types.Yellow)
}
