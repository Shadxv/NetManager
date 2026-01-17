package model

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
)

func CreateNewMongoService(config interfaces.MongoConfig, serviceManager interfaces.ServiceManager) {
	serviceManager.AddService(
		"mongodb",
		types.MongoDB,
		types.Starting,
		"mongodb",
		"mongodb",
		"8.0.4",
		NewMongoData(
			config.GetServiceUsername(),
			config.GetServicePassword(),
			config.GetInternalURI(),
			config.GetExternalURI(),
		),
	)
}
