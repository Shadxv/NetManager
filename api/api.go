package api

import (
	"NetManager/api/controllers"
	"NetManager/internal/module"
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	router        *chi.Mux
	httpServer    *http.Server
	moduleManager *module.Manager
	printer       interfaces.Printer
	status        string
	service       interfaces.Service
	wg            *sync.WaitGroup
}

func NewServer(wg *sync.WaitGroup) *Server {
	s := &Server{
		router: chi.NewRouter(),
		service: interfaces.Service{
			Name: "Gateway",
		},
		wg: wg,
	}
	return s
}

func (s *Server) routes() {
	s.router.Route("/gateway", func(r chi.Router) {
		r.Mount("/v1", controllers.NewV1(s.moduleManager, s.printer, s.service).Router())
	})

}

func (s *Server) Init(moduleManager *module.Manager) {
	s.status = types.Starting
	s.moduleManager = moduleManager
	printer, err := module.GetTypedModule[interfaces.Printer](moduleManager, types.Console)
	if err != nil {
		s.status = types.Disabled
		return
	}
	s.printer = printer

	s.routes()

	s.httpServer = &http.Server{
		Addr:    ":4000",
		Handler: s.router,
	}

	s.wg.Add(1)

	go func() {
		s.status = types.Enabled
		s.printer.PrintColored("Web API starting on :4000", s.printer.Service(), types.Green)
		err = s.httpServer.ListenAndServe()
		if err != nil {
			s.printer.PrintColored(err.Error(), s.printer.Service(), types.Red)
			s.status = types.Disabled
			s.wg.Done()
		}
	}()
}

func (s *Server) Disable(shutdown bool) {
	if s.status == types.Disabled || s.httpServer == nil {
		return
	}

	s.status = types.Stopping
	s.printer.PrintColored("Stopping Web API...", s.service, types.Yellow)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.printer.PrintColored("API Force Shutdown: "+err.Error(), s.service, types.Red)
	}

	s.status = types.Disabled
	s.printer.Print("Web API Disabled", s.service)
	s.wg.Done()
}

func (s *Server) Reload() {
	s.Disable(false)
	s.Init(s.moduleManager)
}

func (s *Server) SaveData() error {
	return nil
}

func (s *Server) LoadData() {}

func (s *Server) SetStatus(newStatus string) {
	s.status = newStatus
}

func (s *Server) Status() string {
	return s.status
}

func (s *Server) Type() string {
	return types.Gateway
}
