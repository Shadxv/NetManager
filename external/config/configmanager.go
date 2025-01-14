package config

import configModel "NetManager/internal/config/model"

type Manager interface {
	Init()
	GetMainConfig() *configModel.MainConfig
	GetRedisConfig() *configModel.RedisConfig
	GetHarborConfig() *configModel.HarborConfig
}
