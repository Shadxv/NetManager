package interfaces

type VelocityData interface {
	Version() string
	BuildNumber() int
	Port() int
	ReplicasAmount() int
}

func GetVelocityData(service ServiceModel) VelocityData {
	data := *service.ServiceData()
	if data, ok := data.(VelocityData); ok {
		return data
	}
	return nil
}
