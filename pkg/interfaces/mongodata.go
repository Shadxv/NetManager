package interfaces

type MongoData interface {
	Port() int
	ExternalPort() int
	Username() string
	Password() string
	AuthRequired() bool
	Authorization() bool
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
