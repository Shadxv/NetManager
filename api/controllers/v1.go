package controllers

import (
	"NetManager/api/auth"
	v1Controller "NetManager/api/controllers/v1"
	"NetManager/internal/module"
	"NetManager/pkg/interfaces"

	"github.com/go-chi/chi/v5"
)

type V1 struct {
	moduleManager *module.Manager
	printer       interfaces.Printer
	service       interfaces.Service
	r             *chi.Mux
	jwtSecret     string
}

func NewV1(moduleManager *module.Manager, printer interfaces.Printer, service interfaces.Service, jwtSecret string) *V1 {
	return &V1{
		moduleManager: moduleManager,
		printer:       printer,
		service:       service,
		r:             chi.NewRouter(),
		jwtSecret:     jwtSecret,
	}
}

func (v1 *V1) Router() *chi.Mux {
	v1.r.Group(func(r chi.Router) {
		r.Use(auth.AuthenticateToken(v1.jwtSecret, false))

		r.Route("/services", func(r chi.Router) {
			serviceHandler := v1Controller.NewServiceHandler(v1.moduleManager, v1.printer, v1.service)
			serviceHandler.RegisterRoutes(r)
		})
	})

	return v1.r
}
