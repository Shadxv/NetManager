package interfaces

type PaperData interface {
	Version() string
	BuildNumber() int
	MinReplicas() int
}

func GetPaperData(service ServiceModel) *PaperData {
	if data, ok := service.ServiceData().(PaperData); ok {
		return &data
	}
	return nil
}
