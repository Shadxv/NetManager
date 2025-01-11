package model

import (
	"NetManager/external/types"
	"slices"
)

type Service struct {
	name              string
	serviceType       string
	status            string
	imageName         string
	currentVersion    string
	availableVersions []string
	serviceData       *interface{}
}

func NewService(name string, serviceType string, status string, imageName string, currentVersion string, availableVersions []string, serviceData *interface{}) *Service {
	return &Service{
		name:              name,
		serviceType:       serviceType,
		status:            status,
		imageName:         imageName,
		currentVersion:    currentVersion,
		availableVersions: availableVersions,
		serviceData:       serviceData,
	}
}

func CreateNewService(name string, serviceType string, serviceData *interface{}) *Service {
	return &Service{
		name:        name,
		serviceType: serviceType,
		status:      types.Stopped,
		imageName:   name,
		serviceData: serviceData,
	}
}

func (service *Service) Name() string {
	return service.name
}

func (service *Service) ServiceType() string {
	return service.serviceType
}

func (service *Service) Status() string {
	return service.status
}

func (service *Service) SetStatus(status string) {
	service.status = status
}

func (service *Service) ImageName() string {
	return service.imageName
}

func (service *Service) CurrentVersion() string {
	return service.currentVersion
}

func (service *Service) SetCurrentVersion(currentVersion string) bool {
	if slices.Contains(service.availableVersions, currentVersion) {
		service.currentVersion = currentVersion
		return true
	}
	return false
}

func (service *Service) AvailableVersions() []string {
	versionsCopy := make([]string, len(service.availableVersions))
	copy(versionsCopy, service.availableVersions)
	return versionsCopy
}

func (service *Service) AddVersion(version string) {
	service.availableVersions = append(service.availableVersions, version)
}

func (service *Service) ServiceData() *interface{} {
	return service.serviceData
}
