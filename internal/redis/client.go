package redis

import (
	"NetManager/internal/redis/handler"
	"NetManager/internal/redis/model"
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"strings"
)

type Client struct {
	printer        interfaces.Printer
	clusterManager interfaces.ClusterManager
	data           interfaces.RedisData
	client         *redis.Client
	codec          interfaces.PacketCodec
	listeners      map[string]interfaces.PacketHandler
	pubsub         *redis.PubSub
}

func NewClient(printer interfaces.Printer, clusterManager interfaces.ClusterManager) *Client {
	return &Client{
		printer:        printer,
		clusterManager: clusterManager,
		codec:          model.NewPacketCodec(),
		listeners:      map[string]interfaces.PacketHandler{},
	}
}

func (redisClient *Client) Init(service interfaces.ServiceModel) {
	redisClient.printer.Print("Connecting to Redis client...", redisClient.printer.Service())
	redisClient.data = interfaces.GetRedisData(service)

	if redisClient.data == nil {
		redisClient.printer.PrintColored("Redis data not found.", redisClient.printer.Service(), types.Red)
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

	redisClient.printer.Print("Successfully connected to Redis client...", redisClient.printer.Service())
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
}

func (redisClient *Client) BuildChannel(groupName string, serviceName string, serviceId string, packetType string) string {
	return fmt.Sprintf("%s:%s:%s:%s", groupName, serviceName, serviceId, packetType)
}

func (redisClient *Client) Close() {
	redisClient.pubsub.Close()
	redisClient.client.Close()
}

func (redisClient *Client) listen() {
	for msg := range redisClient.pubsub.Channel() {
		channel := msg.Channel
		packetType := strings.Split(channel, ":")[3]

		handler, ok := redisClient.listeners[packetType]
		packet, err := redisClient.codec.UnmarshalPacket(packetType, []byte(msg.Payload))

		if !ok || err != nil {
			redisClient.printer.PrintColored("Unable to handle "+packetType+" packet.", redisClient.printer.Service(), types.Red)
			continue
		}

		handler.Handle(packet)
	}
}

func (redisClient *Client) registerListeners() {
	handler.NewRegisterServerRequestListener(redisClient, redisClient.clusterManager)
}
