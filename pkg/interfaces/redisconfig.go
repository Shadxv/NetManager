package interfaces

type RedisConfig interface {
	GetDockerImage() string
	GetVersion() string
	GetPort() int
	GetPassword() string
	GetMaxMemory() string
	GetMaxMemoryPolicy() string
}
