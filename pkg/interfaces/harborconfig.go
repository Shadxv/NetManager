package interfaces

type HarborConfig interface {
	GetDomain() string
	GetProjectName() string
	GetUsername() string
	GetPassword() string
}
