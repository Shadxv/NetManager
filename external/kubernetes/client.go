package kubernetes

type Client interface {
	Connect()
	GetNamespace()
	DeployRedis()
	DeployPaperService(serviceName string)
}
