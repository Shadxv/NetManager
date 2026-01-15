package interfaces

type PaperData interface {
	GroupName() string
	Version() string
	BuildNumber() int
	MinReplicas() int
	MongoDBURI() string
	RedisURI() string
}

func GetPaperData(service ServiceModel) PaperData {
	data := service.ServiceData()
	if data, ok := data.(PaperData); ok {
		return data
	}
	return nil
}
