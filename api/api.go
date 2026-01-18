package api

import (
	"NetManager/api/controllers"
	"NetManager/api/sse"
	"NetManager/internal/module"
	serviceManager "NetManager/internal/service/manager"
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
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
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

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

	svcManager, err := module.GetTypedModule[*serviceManager.ServiceManager](moduleManager, types.Services)
	if err == nil {
		sseServer := sse.NewServer()
		sseServer.SetServiceManager(svcManager)
		go sseServer.Run()
		svcManager.SetBroadcaster(sseServer)
		s.router.Handle("/gateway/events/service", sseServer)
		s.router.Handle("/gateway/events/pod", sseServer)
	} else {
		s.printer.PrintColored("Could not initialize SSE: ServiceManager not found", s.service, types.Red)
	}

	s.httpServer = &http.Server{
		Addr:    ":4000",
		Handler: s.router,
	}

	s.wg.Add(1)

	go func() {
		s.status = types.Enabled
		s.printer.PrintColored("Web API starting on :4000", s.printer.Service(), types.Green)
		err = s.httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
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
	s.printer.Print("Stopping Web API...", s.service)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		s.printer.Print("API Force Shutdown: "+err.Error(), s.service)
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
