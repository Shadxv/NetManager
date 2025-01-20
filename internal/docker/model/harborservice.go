package model

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
)

func CreateHarborService(printer interfaces.Printer, config interfaces.HarborConfig, serviceManager interfaces.ServiceManager, clusterManager interfaces.ClusterManager) {
	serviceModel := serviceManager.AddService(
		"harbor",
		types.Harbor,
		types.Running,
		"goharbor",
		"",
		NewHarborData(
			config.GetDomain(),
			config.GetProjectName(),
			config.GetUsername(),
			config.GetPassword(),
		),
	)

	serviceModel.Deploy(printer, clusterManager)
}
