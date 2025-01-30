package model

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
)

func CreateNewMongoService(printer interfaces.Printer, config interfaces.MongoConfig, serviceManager interfaces.ServiceManager, clusterManager interfaces.ClusterManager) {
	serviceModel := serviceManager.AddService(
		"mongodb",
		types.MongoDB,
		types.Starting,
		"mongodb",
		"latest",
		NewMongoData(
			config.GetPort(),
			config.GetExternalPort(),
			config.GetRootUsername(),
			config.GetRootPassword(),
			config.GetServiceUsername(),
			config.GetServicePassword(),
			config.IsAuthRequired(),
		),
	)

	serviceModel.Deploy(printer, clusterManager)
}
