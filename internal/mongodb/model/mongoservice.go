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
			config.GetUsername(),
			config.GetPassword(),
			config.IsAuthRequired(),
			config.NeedsAuthorization(),
		),
	)

	serviceModel.Deploy(printer, clusterManager)
}
