package docker

import "NetManager/external/service"

type HarborData interface {
	Domain() string
	ProjectName() string
	Username() string
	Password() string
}

func GetHarborData(service service.Model) HarborData {
	serviceData := *service.ServiceData()
	if data, ok := serviceData.(HarborData); ok {
		return data
	}
	return nil
}
