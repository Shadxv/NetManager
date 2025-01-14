package kubernetes

import "NetManager/internal/kubernetes/manager"

type Client interface {
	Connect()
	IsLoaded() bool
	ClusterManager() *manager.ClusterManager
}
