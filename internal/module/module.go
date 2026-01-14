package module

type Module interface {
	Init(moduleManager *Manager)
	Disable(shutdown bool)
	Reload()
	SaveData() error
	LoadData()
	SetStatus(newStatus string)
	Status() string
	Type() string
}
