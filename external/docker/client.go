package docker

type Client interface {
	Init()
	Close()
}
