package model

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
)

type MongoData struct {
	port            int
	externalPort    int
	username        string
	password        string
	authRequired    bool
	authorization   bool
	internalMongoIp string
	externalMongoIp string
}

func NewMongoData(port int, externalPort int, username string, password string, authRequired bool, authorization bool, internalMongoIp string, externalMongoIp string) *MongoData {
	return &MongoData{
		port:            port,
		externalPort:    externalPort,
		username:        username,
		password:        password,
		authRequired:    authRequired,
		authorization:   authorization,
		internalMongoIp: internalMongoIp,
		externalMongoIp: externalMongoIp,
	}
}

func (data *MongoData) Port() int {
	return data.port
}

func (data *MongoData) ExternalPort() int {
	return data.externalPort
}

func (data *MongoData) Username() string {
	return data.username
}

func (data *MongoData) Password() string {
	return data.password
}

func (data *MongoData) AuthRequired() bool {
	return data.authRequired
}

func (data *MongoData) Authorization() bool {
	return data.authorization
}

func (data *MongoData) InternalMongoIp() string {
	return data.internalMongoIp
}

func (data *MongoData) ExternalMongoIp() string {
	return data.externalMongoIp
}

func (data *MongoData) Build(serviceModel interfaces.ServiceModel, printer interfaces.Printer, imageManager interfaces.ImageManager, serviceManager interfaces.ServiceManager) {
	printer.PrintColored("MongoDB service cannot be built.", printer.Service(), types.Yellow)
}

func (data *MongoData) Update(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	printer.PrintColored("MongoDB service cannot be updated.", printer.Service(), types.Yellow)
}

func (data *MongoData) Stop(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	printer.Print("Stopping MongoDB service...", printer.Service())

	printer.Print("MongoDB service has been stopped", printer.Service())
}

func (data *MongoData) Start(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	printer.PrintColored("MongoDB service cannot be started.", printer.Service(), types.Yellow)
}

func (data *MongoData) Deploy(serviceModel interfaces.ServiceModel, printer interfaces.Printer, clusterManager interfaces.ClusterManager) {
	printer.Print("Checking MongoDB services statuses...", printer.Service())

}
