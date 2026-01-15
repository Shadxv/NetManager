package interfaces

import "NetManager/internal/module"

type RedisClient interface {
	Init(moduleManager *module.Manager)
	Publish(channel string, packet Packet)
	RegisterListener(handler PacketHandler)
	BuildChannel(groupName string, serviceName string, serviceId string, packetType string) string
}
