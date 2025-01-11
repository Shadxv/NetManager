package service

type Manager interface {
	AddNewService(name string, serviceType string, serviceData *interface{}) Model
	DeleteService(name string) Model
	Exists(name string) bool
	GetService(name string) Model
	Services() []Model
}
