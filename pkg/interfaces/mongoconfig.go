package interfaces

type MongoConfig interface {
	GetPort() int
	GetExternalPort() int
	GetUsername() string
	GetPassword() string
	IsAuthRequired() bool
	NeedsAuthorization() bool
}
