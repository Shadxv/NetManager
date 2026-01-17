package module

import (
	"NetManager/pkg/types"
	"NetManager/pkg/util"
	"errors"
	"sync"
)

type Manager struct {
	modules *util.OrderedMap[string, Module]
	wg      *sync.WaitGroup
}

func (manager *Manager) Init() {
	for _, item := range manager.modules.Items() {
		module := item.Value
		module.SetStatus(types.Starting)
		module.Init(manager)
	}
}

func (manager *Manager) Disable() {
	items := manager.modules.Items()
	for i := len(items) - 1; i >= 0; i-- {
		module := items[i].Value
		module.Disable(true)
	}
}

func (manager *Manager) AddModule(module Module) {
	if _, exists := manager.modules.Get(module.Type()); exists {
		return
	}

	manager.modules.Set(module.Type(), module)
}

func (manager *Manager) RemoveModule(module Module, shutdown bool) {
	module.Disable(shutdown)
	manager.modules.Delete(module.Type())
}

func (manager *Manager) GetModules() []Module {
	values := make([]Module, 0, manager.modules.Length())
	for _, v := range manager.modules.Items() {
		values = append(values, v.Value)
	}
	return values
}

func (manager *Manager) GetModule(moduleName string) Module {
	module, exists := manager.modules.Get(moduleName)
	if exists {
		return module
	}
	return nil
}

func GetTypedModule[T any](manager *Manager, moduleName string) (T, error) {
	var zero T
	module := manager.GetModule(moduleName)
	if module == nil || module.Status() != types.Enabled {
		return zero, errors.New("module not found or is not enabled")
	}
	casted, ok := module.(T)
	if !ok {
		return zero, errors.New("module is not T")
	}
	return casted, nil
}

func (manager *Manager) AddToWG() {
	manager.wg.Add(1)
}

func NewModuleManager(mainWaitGroup *sync.WaitGroup) *Manager {
	return &Manager{
		modules: util.NewOrderedMap[string, Module](),
		wg:      mainWaitGroup,
	}
}
