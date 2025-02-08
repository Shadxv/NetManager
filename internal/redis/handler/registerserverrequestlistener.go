package handler

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/packets"
	"NetManager/pkg/types"
	"fmt"
)

type RegisterServerRequestListener struct {
	redisClient    interfaces.RedisClient
	packetType     string
	clusterManager interfaces.ClusterManager
}

func NewRegisterServerRequestListener(redisClient interfaces.RedisClient, clusterManager interfaces.ClusterManager) {
	listener := &RegisterServerRequestListener{
		redisClient:    redisClient,
		packetType:     types.RegisterServerRequest,
		clusterManager: clusterManager,
	}

	redisClient.RegisterListener(listener)
}

func (listener *RegisterServerRequestListener) GetType() string {
	return listener.packetType
}

func (listener *RegisterServerRequestListener) Handle(packet interfaces.Packet) error {
	regPacket, ok := packet.(*packets.RegisterServerRequest)

	if !ok {
		return fmt.Errorf("error occured during casting packet")
	}

	// FIXME: Change that so it could handle multiple replicas from 1 service
	service, err := listener.clusterManager.GetServiceOrErr(regPacket.SenderServiceName + "-service")
	if err != nil {
		return fmt.Errorf("error occured during casting packet")
	}

	ip := service.Spec.ClusterIP
	port := 25565

	senderChannel := listener.redisClient.BuildChannel(regPacket.SenderServiceGroup, regPacket.SenderServiceName, regPacket.SenderServiceId, types.CacheIP)
	proxyChannel := listener.redisClient.BuildChannel(regPacket.ProxyGroupName, regPacket.ProxyServiceName, "*", types.RegisterServerData)

	regDataPacket := packets.RegisterServerData{
		PacketData:  regPacket.PacketData,
		AddressIP:   ip,
		Port:        port,
		ServiceName: regPacket.SenderServiceName,
		ServiceId:   regPacket.SenderServiceId,
	}

	regDataPacket.PacketType = types.RegisterServerData

	cacheIP := packets.CacheIP{
		PacketData: interfaces.NewPacketData(types.CacheIP),
		AddressIP:  ip,
		Port:       port,
	}

	go listener.redisClient.Publish(senderChannel, &cacheIP)
	go listener.redisClient.Publish(proxyChannel, &regDataPacket)

	return nil
}

func (listener *RegisterServerRequestListener) GetRegistryType() interfaces.Packet {
	return &packets.RegisterServerRequest{}
}
