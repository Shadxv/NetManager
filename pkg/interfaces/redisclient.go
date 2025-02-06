package interfaces

type RedisClient interface {
	Init(service ServiceModel)
	Publish(channel string, packet Packet)
	RegisterListener(handler PacketHandler)
	BuildChannel(groupName string, serviceName string, serviceId string, packetType string) string
}
