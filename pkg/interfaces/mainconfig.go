package interfaces

type MainConfig interface {
	GetName() string
	GetServerGroupName() string
	GetJWTSecret() string
}
