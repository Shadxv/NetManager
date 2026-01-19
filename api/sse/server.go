package sse

import (
	"NetManager/pkg/interfaces"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Client struct {
	id               string
	send             chan []byte
	subscriptionType string
	targetID         string
	serviceName      string
}

type Server struct {
	clients        map[string]*Client
	register       chan *Client
	unregister     chan *Client
	broadcast      chan Message
	mutex          sync.RWMutex
	serviceManager interfaces.ServiceManager
}

type Message struct {
	ServiceName string      `json:"serviceName"`
	PodName     string      `json:"podName,omitempty"`
	Type        string      `json:"type"`
	Payload     interface{} `json:"payload"`
}

func NewServer() *Server {
	s := &Server{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan Message),
	}
	go s.heartbeat()
	return s
}

func (s *Server) SetServiceManager(sm interfaces.ServiceManager) {
	s.serviceManager = sm
}

func (s *Server) heartbeat() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		s.broadcast <- Message{Type: "ping"}
	}
}

func (s *Server) Run() {
	for {
		select {
		case client := <-s.register:
			s.mutex.Lock()
			s.clients[client.id] = client
			s.mutex.Unlock()

			if s.serviceManager != nil {
				if client.subscriptionType == "pod" {
					service := s.serviceManager.GetService(client.serviceName)
					if service != nil {
						for _, pod := range service.PodInstances() {
							if pod.Name() == client.targetID {
								logs := pod.Logs()
								start := 0
								if len(logs) > 100 {
									start = len(logs) - 100
								}
								for i := start; i < len(logs); i++ {
									msg := Message{
										ServiceName: client.serviceName,
										PodName:     pod.Name(),
										Type:        "log",
										Payload: map[string]string{
											"line": logs[i],
										},
									}
									s.sendToClient(client, msg)
								}
								statusMsg := Message{
									ServiceName: client.serviceName,
									PodName:     pod.Name(),
									Type:        "pod_update",
									Payload: map[string]string{
										"status": pod.Status(),
									},
								}
								s.sendToClient(client, statusMsg)
								break
							}
						}
					}
				}
			}

		case client := <-s.unregister:
			s.mutex.Lock()
			if _, ok := s.clients[client.id]; ok {
				delete(s.clients, client.id)
				close(client.send)
			}
			s.mutex.Unlock()

		case message := <-s.broadcast:
			s.mutex.RLock()
			for _, client := range s.clients {
				if message.Type == "ping" {
					s.sendToClient(client, message)
					continue
				}

				shouldSend := false

				if client.subscriptionType == "service" {
					if client.targetID == "" || client.targetID == message.ServiceName {
						if message.Type != "log" {
							shouldSend = true
						}
					}
				} else if client.subscriptionType == "pod" {
					if client.serviceName == message.ServiceName && client.targetID == message.PodName {
						shouldSend = true
					}
				}

				if shouldSend {
					s.sendToClient(client, message)
				}
			}
			s.mutex.RUnlock()
		}
	}
}

func (s *Server) sendToClient(client *Client, message Message) {
	data, err := json.Marshal(message)
	if err == nil {
		select {
		case client.send <- data:
		default:
		}
	}
}

func (s *Server) BroadcastLog(serviceName string, podName string, logLine string) {
	s.broadcast <- Message{
		ServiceName: serviceName,
		PodName:     podName,
		Type:        "log",
		Payload: map[string]string{
			"line": logLine,
		},
	}
}

func (s *Server) BroadcastEvent(serviceName string, eventType string, payload interface{}) {
	podName := ""
	if m, ok := payload.(map[string]interface{}); ok {
		if p, ok := m["pod"].(string); ok {
			podName = p
		}
	}

	s.broadcast <- Message{
		ServiceName: serviceName,
		PodName:     podName,
		Type:        eventType,
		Payload:     payload,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	path := r.URL.Path
	var subType, targetID, serviceName string

	if strings.HasSuffix(path, "/service") {
		subType = "service"
		targetID = r.URL.Query().Get("name")
	} else if strings.HasSuffix(path, "/pod") {
		subType = "pod"
		targetID = r.URL.Query().Get("name")
		serviceName = r.URL.Query().Get("service")
	} else {
		http.Error(w, "Invalid endpoint", http.StatusBadRequest)
		return
	}

	clientID := fmt.Sprintf("%s-%d-%d", r.RemoteAddr, time.Now().UnixNano(), r.URL.Port())

	client := &Client{
		id:               clientID,
		send:             make(chan []byte, 2048),
		subscriptionType: subType,
		targetID:         targetID,
		serviceName:      serviceName,
	}

	s.register <- client

	notify := w.(http.Flusher)
	fmt.Fprintf(w, "data: %s\n\n", "connected")
	notify.Flush()

	ctx := r.Context()

	for {
		select {
		case <-ctx.Done():
			s.unregister <- client
			return
		case msg, ok := <-client.send:
			if !ok {
				return
			}
			var m Message
			if err := json.Unmarshal(msg, &m); err == nil && m.Type == "ping" {
				fmt.Fprintf(w, ": ping\n\n")
			} else {
				fmt.Fprintf(w, "data: %s\n\n", msg)
			}
			notify.Flush()
		}
	}
}
