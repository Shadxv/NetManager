package interfaces

type HarborData interface {
	Domain() string
	ProjectName() string
	Username() string
	Password() string
}

func GetHarborData(service ServiceModel) HarborData {
	data := *service.ServiceData()
	if data, ok := data.(HarborData); ok {
		return data
	}
	return nil
}
