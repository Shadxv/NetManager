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

func NewVelocityService(name string) *VelocityService {
	return &VelocityService{
		name: name,
	}
}

func (service *VelocityService) GetType() string {
	return "paper"
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
