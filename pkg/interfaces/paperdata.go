package interfaces

type PaperData interface {
	Version() string
	BuildNumber() int
	MinReplicas() int
}

func GetPaperData(service ServiceModel) PaperData {
	data := *service.ServiceData()
	if data, ok := data.(PaperData); ok {
		return data
	}
	return nil
}
