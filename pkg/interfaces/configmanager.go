package interfaces

type ConfigManager interface {
	Init()
	GetMainConfig() MainConfig
	GetRedisConfig() RedisConfig
	GetHarborConfig() HarborConfig
	GetMongoConfig() MongoConfig
}
