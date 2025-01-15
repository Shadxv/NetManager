package model

import (
	"NetManager/external/kubernetes"
	"NetManager/external/service"
	"NetManager/external/types"
	"NetManager/internal/config/model"
)

func CreateHarborService(config *model.HarborConfig, serviceManager service.Manager, clusterManager kubernetes.ClusterManager) {
	serviceModel := serviceManager.AddService(
		"harbor",
		types.Harbor,
		types.Running,
		"goharbor",
		"",
		NewHarborData(
			config.Domain,
			config.ProjectName,
			config.Username,
			config.Password,
		),
	)

	getPods(serviceModel, clusterManager)
}

func getPods(service service.Model, manager kubernetes.ClusterManager) {
	pods := manager.GetPods("app=harbor")
	if len(pods) != 0 {
		for _, pod := range pods {
			service.AddPodInstance(pod.Name, &pod, string(pod.Status.Phase))
		}
	}
}
