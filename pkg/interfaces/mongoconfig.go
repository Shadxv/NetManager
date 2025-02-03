package interfaces

type MongoConfig interface {
	GetServiceUsername() string
	GetServicePassword() string
	GetInternalURI() string
	GetExternalURI() string
}
