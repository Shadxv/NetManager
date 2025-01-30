package interfaces

type MongoConfig interface {
	GetPort() int
	GetExternalPort() int
	GetRootUsername() string
	GetRootPassword() string
	GetServiceUsername() string
	GetServicePassword() string
	IsAuthRequired() bool
}
