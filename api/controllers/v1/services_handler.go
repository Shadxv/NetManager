package v1

import (
	"NetManager/api/auth"
	"NetManager/internal/kubernetes"
	minecraftModel "NetManager/internal/minecraft/model"
	"NetManager/internal/module"
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	corev1 "k8s.io/api/core/v1"
)

type ServiceHandler struct {
	moduleManager  *module.Manager
	printer        interfaces.Printer
	serviceManager interfaces.ServiceManager
	clusterManager interfaces.ClusterManager
	service        interfaces.Service
}

type ServiceBaseInfo struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

type PodInfo struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	InternalIP string `json:"internalIP"`
	ExternalIP string `json:"externalIP"`
	Port       int32  `json:"port"`
}

type ServiceDetails struct {
	Status            string    `json:"status"`
	CurrentVerion     string    `json:"currentVersion"`
	AvailableVersions []string  `json:"availableVersions"`
	Pods              []PodInfo `json:"pods"`
}

type CreateServiceRequest struct {
	Name string          `json:"name"`
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type VersionUpdate struct {
	Version string `json:"version"`
}

type VelocityData struct {
	Version string `json:"version"`
	Build   int    `json:"build"`
	Port    int    `json:"port"`
}

type PaperData struct {
	Version string `json:"version"`
	Build   int    `json:"build"`
}

type StatsResponse struct {
	TotalServices   int `json:"totalServices"`
	RunningServices int `json:"runningServices"`
	TotalPods       int `json:"totalPods"`
	NonRunningPods  int `json:"nonRunningPods"`
}

func NewServiceHandler(moduleManager *module.Manager, printer interfaces.Printer, service interfaces.Service) *ServiceHandler {
	sh := &ServiceHandler{
		moduleManager: moduleManager,
		printer:       printer,
		service:       service,
	}

	serviceManager, err := module.GetTypedModule[interfaces.ServiceManager](sh.moduleManager, types.Services)
	if err == nil {
		sh.serviceManager = serviceManager
	}

	k8sClient, err := module.GetTypedModule[*kubernetes.Client](sh.moduleManager, types.Kubernetes)
	if err == nil {
		sh.clusterManager = k8sClient.ClusterManager()
	}

	return sh
}

func (sh *ServiceHandler) RegisterRoutes(r chi.Router) {
	r.Get("/stats", sh.handleStats)

	r.With(auth.Authorize(auth.READ_SERVICES)).Get("/", sh.handleListServices)
	r.With(auth.Authorize(auth.READ_SERVICES)).Get("/{name}", sh.handleServiceDetails)
	r.With(auth.Authorize(auth.DELETE_SERVICES)).Delete("/{name}", sh.handleDeleteService)
	r.With(auth.Authorize(auth.MANAGE_SERVICES_STATE)).Post("/{name}/start", sh.handleStartService)
	r.With(auth.Authorize(auth.MANAGE_SERVICES_STATE)).Post("/{name}/stop", sh.handleStopService)
	r.With(auth.Authorize(auth.MANAGE_SERVICES_STATE)).Post("/{name}/restart", sh.handleRestartService)
	r.With(auth.Authorize(auth.MANAGE_SERVICES_STATE)).Post("/{name}/pod/{id}/stop", sh.handleStopPod)
	r.With(auth.Authorize(auth.CREATE_NEW_SERVICES)).Post("/create", sh.handleCreateService)
	r.With(auth.Authorize(auth.UPDATE_SERVICES)).Post("/{name}/build", sh.handleBuildVersion)
	r.With(auth.Authorize(auth.UPDATE_SERVICES)).Patch("/{name}/version", sh.handleChangeVersion)
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

func (sh *ServiceHandler) handleStats(w http.ResponseWriter, r *http.Request) {
	services := sh.serviceManager.GetServices()
	stats := StatsResponse{}

	for _, svc := range services {
		if svc.ServiceType() == types.Paper || svc.ServiceType() == types.Velocity {
			stats.TotalServices++
			if svc.Status() == types.Running {
				stats.RunningServices++
			}

			for _, pod := range svc.PodInstances() {
				stats.TotalPods++
				if pod.Status() != string(corev1.PodRunning) {
					stats.NonRunningPods++
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		sh.printer.PrintColored("API JSON Error: "+err.Error(), sh.service, types.Red)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (sh *ServiceHandler) handleServiceDetails(w http.ResponseWriter, r *http.Request) {
	serviceName := chi.URLParam(r, "name")
	service := sh.serviceManager.GetService(serviceName)
	if service == nil || service.ServiceType() != types.Paper && service.ServiceType() != types.Velocity {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	pods := make([]PodInfo, 0)
	for _, pod := range service.PodInstances() {
		podInfo := PodInfo{
			Name:       pod.Name(),
			Status:     pod.Status(),
			InternalIP: pod.InternalIP(),
			ExternalIP: pod.ExternalIP(),
			Port:       pod.Port(),
		}
		pods = append(pods, podInfo)
	}

	details := ServiceDetails{
		Status:            service.Status(),
		CurrentVerion:     service.CurrentVersion(),
		AvailableVersions: service.AvailableVersions(),
		Pods:              pods,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(details); err != nil {
		sh.printer.PrintColored("API JSON Error: "+err.Error(), sh.service, types.Red)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (sh *ServiceHandler) handleDeleteService(w http.ResponseWriter, r *http.Request) {
	serviceName := chi.URLParam(r, "name")
	service := sh.serviceManager.GetService(serviceName)
	if service == nil || service.ServiceType() != types.Paper && service.ServiceType() != types.Velocity {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}
	sh.serviceManager.DeleteService(serviceName)
	w.WriteHeader(http.StatusOK)
}

func (sh *ServiceHandler) handleStartService(w http.ResponseWriter, r *http.Request) {
	serviceName := chi.URLParam(r, "name")
	service := sh.serviceManager.GetService(serviceName)
	if service == nil || (service.ServiceType() != types.Paper && service.ServiceType() != types.Velocity) {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	if sh.clusterManager == nil {
		http.Error(w, "Cluster Manager not available", http.StatusInternalServerError)
		return
	}

	service.Deploy(sh.printer, sh.clusterManager)
	w.WriteHeader(http.StatusOK)
}

func (sh *ServiceHandler) handleStopService(w http.ResponseWriter, r *http.Request) {
	serviceName := chi.URLParam(r, "name")
	service := sh.serviceManager.GetService(serviceName)
	if service == nil || (service.ServiceType() != types.Paper && service.ServiceType() != types.Velocity) {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	if sh.clusterManager == nil {
		http.Error(w, "Cluster Manager not available", http.StatusInternalServerError)
		return
	}

	service.Stop(sh.printer, sh.clusterManager)
	w.WriteHeader(http.StatusOK)
}

func (sh *ServiceHandler) handleRestartService(w http.ResponseWriter, r *http.Request) {
	serviceName := chi.URLParam(r, "name")
	service := sh.serviceManager.GetService(serviceName)
	if service == nil || (service.ServiceType() != types.Paper && service.ServiceType() != types.Velocity) {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	if sh.clusterManager == nil {
		http.Error(w, "Cluster Manager not available", http.StatusInternalServerError)
		return
	}

	service.Stop(sh.printer, sh.clusterManager)
	service.Deploy(sh.printer, sh.clusterManager)
	w.WriteHeader(http.StatusOK)
}

func (sh *ServiceHandler) handleStopPod(w http.ResponseWriter, r *http.Request) {
	serviceName := chi.URLParam(r, "name")
	podID := chi.URLParam(r, "id")

	service := sh.serviceManager.GetService(serviceName)
	if service == nil || (service.ServiceType() != types.Paper && service.ServiceType() != types.Velocity) {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	if sh.clusterManager == nil {
		http.Error(w, "Cluster Manager not available", http.StatusInternalServerError)
		return
	}

	found := false
	for _, pod := range service.PodInstances() {
		if pod.Name() == podID {
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Pod not found in service", http.StatusNotFound)
		return
	}

	err := sh.clusterManager.DeletePod(podID, sh.clusterManager.GetDefaultNamespace())
	if err != nil {
		sh.printer.PrintColored("Error deleting pod: "+err.Error(), sh.service, types.Red)
		http.Error(w, "Failed to delete pod", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (sh *ServiceHandler) handleCreateService(w http.ResponseWriter, r *http.Request) {
	var req CreateServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if sh.serviceManager.Exists(req.Name) {
		http.Error(w, "Service already exists", http.StatusConflict)
		return
	}

	var serviceData interface{}

	mongoURI, redisURI := sh.getURIs()

	switch req.Type {
	case types.Paper:
		var paperData PaperData
		var data minecraftModel.PaperData
		if err := json.Unmarshal(req.Data, &paperData); err != nil {
			http.Error(w, "Invalid data for Paper service", http.StatusBadRequest)
			return
		}
		data = minecraftModel.PaperData{
			GroupNameField:   req.Name,
			VersionField:     paperData.Version,
			BuildField:       paperData.Build,
			MinReplicasField: 1,
			MongodbURIField:  mongoURI,
			RedisURIField:    redisURI,
		}
		serviceData = &data
	case types.Velocity:
		var velocityData VelocityData
		var data minecraftModel.VelocityData
		if err := json.Unmarshal(req.Data, &velocityData); err != nil {
			http.Error(w, "Invalid data for Velocity service", http.StatusBadRequest)
			return
		}
		if velocityData.Port < 30000 || velocityData.Port > 32767 {
			http.Error(w, "Invalid port for Velocity service", http.StatusBadRequest)
		}
		data = minecraftModel.VelocityData{
			GroupNameField:      req.Name,
			VersionField:        velocityData.Version,
			BuildField:          velocityData.Build,
			PortField:           velocityData.Port,
			ReplicasAmountField: 1,
			MongodbURIField:     mongoURI,
			RedisURIField:       redisURI,
		}
		serviceData = &data
	default:
		http.Error(w, "Unsupported service type", http.StatusBadRequest)
		return
	}

	serviceModel := sh.serviceManager.AddNewService(req.Name, req.Type, serviceData)
	result := ServiceBaseInfo{
		Name:   serviceModel.Name(),
		Type:   serviceModel.ServiceType(),
		Status: serviceModel.Status(),
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		sh.printer.PrintColored("API JSON Error: "+err.Error(), sh.service, types.Red)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (sh *ServiceHandler) handleBuildVersion(w http.ResponseWriter, r *http.Request) {
	serviceName := chi.URLParam(r, "name")
	service := sh.serviceManager.GetService(serviceName)
	if service == nil || (service.ServiceType() != types.Paper && service.ServiceType() != types.Velocity) {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	imageManager, err := module.GetTypedModule[interfaces.ImageManager](sh.moduleManager, types.Images)
	if err != nil {
		http.Error(w, "Image Manager not available", http.StatusInternalServerError)
		return
	}

	go service.Build(sh.printer, imageManager, sh.serviceManager)
	w.WriteHeader(http.StatusAccepted)
}

func (sh *ServiceHandler) handleChangeVersion(w http.ResponseWriter, r *http.Request) {
	serviceName := chi.URLParam(r, "name")
	var body VersionUpdate
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	service := sh.serviceManager.GetService(serviceName)

	if service == nil || (service.ServiceType() != types.Paper && service.ServiceType() != types.Velocity) {
		http.Error(w, "Service not found", http.StatusNotFound)
		return
	}

	if !service.SetCurrentVersion(body.Version) {
		http.Error(w, "Invalid version", http.StatusBadRequest)
		return
	}
}

func (sh *ServiceHandler) getURIs() (string, string) {
	mongoService := sh.serviceManager.GetService("mongodb")
	if mongoService == nil {
		return "", ""
	}

	mongoData := interfaces.GetMongoData(mongoService)
	if mongoData == nil {
		return "", ""
	}

	redisService := sh.serviceManager.GetService("redis")
	if redisService == nil {
		return "", ""
	}

	redisData := interfaces.GetRedisData(redisService)
	if redisData == nil {
		return "", ""
	}

	return mongoData.InternalURI(), fmt.Sprintf("redis://:%s@%s:%d/", redisData.Password(), redisData.InternalRedisIp(), redisData.Port())
}
