package interfaces

type VelocityData interface {
	GroupName() string
	Version() string
	BuildNumber() int
	Port() int
	ReplicasAmount() int
	MongoDBURI() string
	RedisURI() string
}

func GetVelocityData(service ServiceModel) VelocityData {
	data := *service.ServiceData()
	if data, ok := data.(VelocityData); ok {
		return data
	}
	return nil
}
