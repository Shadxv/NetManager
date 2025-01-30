package interfaces

type MongoData interface {
	Port() int
	ExternalPort() int
	RootUsername() string
	RootPassword() string
	ServiceUsername() string
	ServicePassword() string
	AuthRequired() bool
	InternalMongoIp() string
	ExternalMongoIp() string
}

func GetMongoData(service ServiceModel) MongoData {
	data := *service.ServiceData()
	if data, ok := data.(MongoData); ok {
		return data
	}
	return nil
}
