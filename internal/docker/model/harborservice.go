package model

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
)

func CreateHarborService(config interfaces.HarborConfig, serviceManager interfaces.ServiceManager, clusterManager interfaces.ClusterManager) {
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

	getPods(serviceModel, clusterManager)
}

func getPods(service interfaces.ServiceModel, manager interfaces.ClusterManager) {
	pods := manager.GetPods("app=harbor")
	if len(pods) != 0 {
		for _, pod := range pods {
			service.AddPodInstance(pod.Name, &pod, string(pod.Status.Phase))
		}
	}
}
