package interfaces

type Data interface {
	Build(serviceModel ServiceModel, printer Printer)
	Update(serviceModel ServiceModel, printer Printer)
	Stop(serviceModel ServiceModel, printer Printer)
	Start(serviceModel ServiceModel, printer Printer)
	Deploy(serviceModel ServiceModel, printer Printer)
}
