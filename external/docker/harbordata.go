package docker

import "NetManager/external/service"

type HarborData interface {
	HttpPort() int
	ProjectName() string
	Username() string
	UserMail() string
	UserPassword() string
	UserRole() string
	DisableGuest() bool
}

func GetHarborData(service service.Model) *HarborData {
	serviceData := *service.ServiceData()
	if data, ok := serviceData.(HarborData); ok {
		return &data
	}
	return nil
}
