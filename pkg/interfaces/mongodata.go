package interfaces

type MongoData interface {
	ServiceUsername() string
	ServicePassword() string
	InternalURI() string
	ExternalURI() string
}

func GetMongoData(service ServiceModel) MongoData {
	data := *service.ServiceData()
	if data, ok := data.(MongoData); ok {
		return data
	}
	return nil
}
