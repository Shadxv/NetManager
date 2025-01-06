package types

import serviceModel "NetManager/internal/minecraft/model"

type ServiceManagerBase interface {
	GetAllServices() map[string]serviceModel.Service
	RegisterNewService(service serviceModel.Service)
}
