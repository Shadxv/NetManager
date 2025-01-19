package interfaces

type HarborData interface {
	Domain() string
	ProjectName() string
	Username() string
	Password() string
}

func GetHarborData(service ServiceModel) HarborData {
	if data, ok := service.ServiceData().(HarborData); ok {
		return data
	}
	return nil
}
