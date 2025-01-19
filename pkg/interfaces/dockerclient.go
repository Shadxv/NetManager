package interfaces

type DockerClient interface {
	Init()
	Close()
}
