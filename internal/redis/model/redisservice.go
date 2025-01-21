package model

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
)

func CreateNewRedisService(printer interfaces.Printer, config interfaces.RedisConfig, serviceManager interfaces.ServiceManager, clusterManager interfaces.ClusterManager) {
	serviceModel := serviceManager.AddService(
		"redis",
		types.Redis,
		types.Starting,
		config.GetDockerImage(),
		config.GetVersion(),
		NewRedisData(
			config.GetPort(),
			config.GetExternalPort(),
			config.GetPassword(),
			config.GetMaxMemory(),
			config.GetMaxMemoryPolicy(),
		),
	)

	serviceModel.Deploy(printer, clusterManager)
}
