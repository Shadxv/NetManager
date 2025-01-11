package minecraft

import "NetManager/external/service"

type VelocityData interface {
	Version() string
	Build() int
	Port() int
	ReplicasAmount() int
}

func GetVelocityData(service service.Model) *VelocityData {
	serviceData := *service.ServiceData()
	if data, ok := serviceData.(VelocityData); ok {
		return &data
	}
	return nil
}
