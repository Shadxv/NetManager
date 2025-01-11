package docker

type ImageManager interface {
	Init()
	BuildImage(serviceName string)
}
