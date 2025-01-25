package interfaces

type MongoData interface {
	Port() int
	ExternalPort() int
	Username() string
	Password() string
	AuthRequired() bool
	Authorization() bool
	InternalMongoIp() string
	ExternalMongoIp() string
}
