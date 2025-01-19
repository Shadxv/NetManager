package interfaces

type VelocityData interface {
	Version() string
	BuildNumber() int
	Port() int
	ReplicasAmount() int
}

func GetVelocityData(service ServiceModel) *VelocityData {
	if data, ok := service.ServiceData().(VelocityData); ok {
		return &data
	}
	return nil
}
