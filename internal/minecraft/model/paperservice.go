package model

type PaperServiceImpl interface {
	GetMinReplicasAmount() int
}

type PaperService struct {
	name              string
	version           string
	build             int
	minReplicasAmount int
}

func NewPaperService(name string, version string, build int, minReplicasAmount int) *PaperService {
	return &PaperService{
		name:              name,
		version:           version,
		build:             build,
		minReplicasAmount: minReplicasAmount,
	}
}

func (service *PaperService) GetType() string {
	return PaperType
}

func (service *PaperService) GetName() string {
	return service.name
}

func (service *PaperService) GetVersion() string {
	return service.version
}

func (service *PaperService) GetBuild() int {
	return service.build
}

func (service *PaperService) GetMinReplicasAmount() int {
	return service.minReplicasAmount
}
