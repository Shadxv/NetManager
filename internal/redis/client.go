package redis

import (
	"NetManager/internal/kubernetes"
	"NetManager/internal/module"
	"NetManager/internal/redis/handler"
	"NetManager/internal/redis/model"
	serviceManager "NetManager/internal/service/manager"
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	printer        interfaces.Printer
	clusterManager interfaces.ClusterManager
	data           interfaces.RedisData
	client         *redis.Client
	codec          interfaces.PacketCodec
	listeners      map[string]interfaces.PacketHandler
	pubsub         *redis.PubSub
	status         string
}

func NewRedisClient() *Client {
	return &Client{
		codec:     model.NewPacketCodec(),
		listeners: map[string]interfaces.PacketHandler{},
		status:    types.Stopped,
	}
}

func (redisClient *Client) Init(moduleManager *module.Manager) {
	redisClient.status = types.Starting

	printer, err := module.GetTypedModule[interfaces.Printer](moduleManager, types.Console)
	if err != nil {
		redisClient.status = types.Disabled
		return
	}
	redisClient.printer = printer

	k8sClient, err := module.GetTypedModule[*kubernetes.Client](moduleManager, types.Kubernetes)
	if err != nil {
		redisClient.printer.PrintColored("Kubernetes module not found!", redisClient.printer.Service(), types.Red)
		redisClient.status = types.Disabled
		return
	}
	redisClient.clusterManager = k8sClient.ClusterManager()

	svcManager, err := module.GetTypedModule[*serviceManager.ServiceManager](moduleManager, types.Services)
	if err != nil {
		redisClient.printer.PrintColored("ServiceManager module not found!", redisClient.printer.Service(), types.Red)
		redisClient.status = types.Disabled
		return
	}

	service := svcManager.GetService("redis")
	if service == nil {
		redisClient.printer.PrintColored("Redis service not found!", redisClient.printer.Service(), types.Red)
		redisClient.status = types.Disabled
		return
	}

	redisClient.printer.Print("Connecting to Redis client...", redisClient.printer.Service())
	redisClient.data = interfaces.GetRedisData(service)

	if redisClient.data == nil {
		redisClient.printer.PrintColored("Redis data not found.", redisClient.printer.Service(), types.Red)
		redisClient.status = types.Disabled
		return
	}

	redisClient.client = redis.NewClient(&redis.Options{
		Addr:     redisClient.data.ExternalRedisIp() + ":" + strconv.Itoa(redisClient.data.ExternalPort()),
		Password: redisClient.data.Password(),
		DB:       0,
		Protocol: 2,
	})

	redisClient.registerListeners()

	var channels []string
	for key := range redisClient.listeners {
		channels = append(channels, redisClient.BuildChannel("netmanager", "netmanager", "*", key), redisClient.BuildChannel("netmanager", "*", "*", key))
	}

	redisClient.pubsub = redisClient.client.Subscribe(context.Background(), channels...)
	go redisClient.listen()

	redisClient.printer.Print("Successfully connected to Redis client.", redisClient.printer.Service())
	redisClient.status = types.Enabled
}

func (redisClient *Client) Disable(shutdown bool) {
	if !shutdown {
		redisClient.printer.PrintColored("This module cannot be disabled!", redisClient.printer.Service(), types.Red)
		return
	}
	redisClient.status = types.Stopping
	redisClient.Close()
	redisClient.status = types.Disabled
}

func (redisClient *Client) Reload() {
	redisClient.printer.PrintColored("This module cannot be reloaded!", redisClient.printer.Service(), types.Red)
}

func (redisClient *Client) SaveData() error {
	return nil
}

func (redisClient *Client) LoadData() {
}

func (redisClient *Client) SetStatus(newStatus string) {
	redisClient.status = newStatus
}

func (redisClient *Client) Status() string {
	return redisClient.status
}

func (redisClient *Client) Type() string {
	return types.RedisClient
}

func (redisClient *Client) Publish(channel string, packet interfaces.Packet) {
	payload, err := redisClient.codec.MarshalPacket(packet)

	if err != nil {
		redisClient.printer.PrintColored("Error occured during payload serialization.", redisClient.printer.Service(), types.Red)
		return
	}

	redisClient.client.Publish(context.Background(), channel, payload)
}

func (redisClient *Client) RegisterListener(handler interfaces.PacketHandler) {
	redisClient.listeners[handler.GetType()] = handler
	redisClient.codec.RegisterType(handler.GetType(), handler.GetRegistryType)
}

func (redisClient *Client) BuildChannel(groupName string, serviceName string, serviceId string, packetType string) string {
	return fmt.Sprintf("%s:%s:%s:%s", groupName, serviceName, serviceId, packetType)
}

func (redisClient *Client) Close() {
	if redisClient.pubsub != nil {
		redisClient.pubsub.Close()
	}
	if redisClient.client != nil {
		redisClient.client.Close()
	}
}

func (redisClient *Client) listen() {
	for msg := range redisClient.pubsub.Channel() {
		channel := msg.Channel
		packetType := strings.Split(channel, ":")[3]

		handler, ok := redisClient.listeners[packetType]
		packet, err := redisClient.codec.UnmarshalPacket(packetType, []byte(msg.Payload))

		if !ok || err != nil {
			redisClient.printer.PrintColored("Unable to handle "+packetType+" packet.", redisClient.printer.Service(), types.Red)

			redisClient.printer.PrintColored(msg.Payload, redisClient.printer.Service(), types.Yellow)

			continue
		}

		handler.Handle(packet)
	}
}

func (redisClient *Client) registerListeners() {
	handler.NewRegisterServerRequestListener(redisClient, redisClient.clusterManager)
}
