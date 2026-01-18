package interfaces

type Broadcaster interface {
	BroadcastLog(serviceName string, podName string, logLine string)
	BroadcastEvent(serviceName string, eventType string, payload interface{})
}
