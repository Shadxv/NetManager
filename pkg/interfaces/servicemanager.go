package interfaces

type ServiceManager interface {
	AddNewService(name string, serviceType string, serviceData *interface{}) ServiceModel
	AddService(name string, serviceType string, status string, image string, version string, serviceData interface{}) ServiceModel
	DeleteService(name string) ServiceModel
	Exists(name string) bool
	GetService(name string) ServiceModel
	Services() []ServiceModel
}
