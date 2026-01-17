package v1

import (
	"NetManager/internal/module"
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ServiceHandler struct {
	moduleManager  *module.Manager
	printer        interfaces.Printer
	serviceManager interfaces.ServiceManager
	r              *chi.Mux
	service        interfaces.Service
}

type ServiceBaseInfo struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

func NewServiceHandler(moduleManager *module.Manager, printer interfaces.Printer, service interfaces.Service) *ServiceHandler {
	return &ServiceHandler{
		moduleManager: moduleManager,
		printer:       printer,
		r:             chi.NewRouter(),
		service:       service,
	}
}

func (sh *ServiceHandler) Router() *chi.Mux {
	serviceManager, err := module.GetTypedModule[interfaces.ServiceManager](sh.moduleManager, types.Services)
	if err != nil {
		return sh.r
	}
	sh.serviceManager = serviceManager

	sh.r.Get("/", sh.handleListServices)

	return sh.r
}

func (sh *ServiceHandler) handleListServices(w http.ResponseWriter, r *http.Request) {
	services := sh.serviceManager.GetServices()

	response := make([]ServiceBaseInfo, 0)

	for _, svc := range services {
		if svc.ServiceType() == types.Paper || svc.ServiceType() == types.Velocity {
			baseInfo := ServiceBaseInfo{
				Name:   svc.Name(),
				Type:   svc.ServiceType(),
				Status: svc.Status(),
			}
			response = append(response, baseInfo)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		sh.printer.PrintColored("API JSON Error: "+err.Error(), sh.service, types.Red)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
