package interfaces

type RedisData interface {
	Port() int
	ExternalPort() int
	Password() string
	MaxMemory() string
	MaxMemoryPolicy() string
	InternalRedisIp() string
	ExternalRedisIp() string
}

func GetRedisData(service ServiceModel) RedisData {
	data := *service.ServiceData()
	if data, ok := data.(RedisData); ok {
		return data
	}
	return nil
}
