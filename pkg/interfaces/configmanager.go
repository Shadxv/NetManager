package interfaces

import "NetManager/internal/module"

type ConfigManager interface {
	Init(moduleManager *module.Manager)
	GetMainConfig() MainConfig
	GetRedisConfig() RedisConfig
	GetHarborConfig() HarborConfig
	GetMongoConfig() MongoConfig
}
