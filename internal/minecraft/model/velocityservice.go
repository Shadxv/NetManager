package model

type VelocityServiceImpl interface {
	GetPort() int
}

type VelocityService struct {
	name    string
	version string
	build   int
	port    int
}

func NewVelocityService(name string, version string, build int, port int) *VelocityService {
	return &VelocityService{
		name:    name,
		version: version,
		build:   build,
		port:    port,
	}
}

func (service *VelocityService) GetType() string {
	return VelocityType
}

func (service *VelocityService) GetName() string {
	return service.name
}

func (service *VelocityService) GetVersion() string {
	return service.version
}

func (service *VelocityService) GetBuild() int {
	return service.build
}

func (service *VelocityService) GetPort() int {
	return service.port
}
