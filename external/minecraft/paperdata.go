package minecraft

import "NetManager/external/service"

type PaperData interface {
	Version() string
	Build() int
	MinReplicas() int
}

func GetPaperData(service service.Model) *PaperData {
	serviceData := *service.ServiceData()
	if data, ok := serviceData.(PaperData); ok {
		return &data
	}
	return nil
}
