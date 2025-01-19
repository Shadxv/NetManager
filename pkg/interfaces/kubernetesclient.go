package interfaces

type KubernetesClient interface {
	Connect()
	IsLoaded() bool
	ClusterManager() ClusterManager
}
