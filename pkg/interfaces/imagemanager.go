package interfaces

import "NetManager/internal/module"

type ImageManager interface {
	Init(manager *module.Manager)
	Client() DockerClient
	BuildImage(serviceModel ServiceModel) (bool, string)
	TagImage(serviceModel ServiceModel, serviceManager ServiceManager) (bool, HarborData)
	PushImage(serviceModel ServiceModel, harborData HarborData) bool
	RemoveImage(id string)
	FullDeployImage(serviceModel ServiceModel, serviceManager ServiceManager) bool
}
