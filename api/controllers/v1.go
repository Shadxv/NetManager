package controllers

import (
	"NetManager/internal/module"
	"NetManager/pkg/interfaces"

	v1Controller "NetManager/api/controllers/v1"

	"github.com/go-chi/chi/v5"
)

type V1 struct {
	moduleManager *module.Manager
	printer       interfaces.Printer
	service       interfaces.Service
	r             *chi.Mux
}

func NewV1(moduleManager *module.Manager, printer interfaces.Printer, service interfaces.Service) *V1 {
	return &V1{
		moduleManager: moduleManager,
		printer:       printer,
		service:       service,
		r:             chi.NewRouter(),
	}
}

func (v1 *V1) Router() *chi.Mux {
	v1.r.Mount("/services", v1Controller.NewServiceHandler(v1.moduleManager, v1.printer, v1.service).Router())

	return v1.r
}
