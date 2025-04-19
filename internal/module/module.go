package module

type Data struct {
}

type Object interface {
	Init()
	Disable()
	Reload()
	SaveData()
	LoadData()
}

const (
	Console    = "Console"
	Config     = "Config"
	Redis      = "Redis"
	MongoDB    = "MongoDB"
	Images     = "Images"
	Kubernetes = "Kubernetes"
)
