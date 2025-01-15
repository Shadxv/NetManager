package docker

import "NetManager/external/service"

type ImageManager interface {
	Init()
	Client() Client
	BuildImage(serviceModel service.Model) (bool, string)
	TagImage(serviceModel service.Model, serviceManager service.Manager) (bool, HarborData)
	PushImage(serviceModel service.Model, harborData HarborData) bool
	RemoveImage(id string)
	FullDeployImage(serviceModel service.Model, serviceManager service.Manager) bool
}
