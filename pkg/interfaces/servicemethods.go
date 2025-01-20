package interfaces

type Data interface {
	Build(serviceModel ServiceModel, printer Printer, imageManager ImageManager, serviceManager ServiceManager)
	Update(serviceModel ServiceModel, printer Printer, clusterManager ClusterManager)
	Stop(serviceModel ServiceModel, printer Printer, clusterManager ClusterManager)
	Start(serviceModel ServiceModel, printer Printer, clusterManager ClusterManager)
	Deploy(serviceModel ServiceModel, printer Printer, clusterManager ClusterManager)
}
