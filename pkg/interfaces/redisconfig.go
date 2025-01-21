package interfaces

type RedisConfig interface {
	GetDockerImage() string
	GetVersion() string
	GetPort() int
	GetExternalPort() int
	GetPassword() string
	GetMaxMemory() string
	GetMaxMemoryPolicy() string
}
