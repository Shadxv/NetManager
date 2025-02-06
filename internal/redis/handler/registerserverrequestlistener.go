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

	service, err := listener.clusterManager.GetServiceOrErr(packet.GetData().SenderServiceName)
	if !ok || err != nil {
		return fmt.Errorf("error occured during casting packet")
	}

	ip := service.Spec.ClusterIP
	port := 25565

	senderChannel := listener.redisClient.BuildChannel(regPacket.Data.SenderServiceGroup, regPacket.Data.SenderServiceName, regPacket.Data.SenderServiceId, types.CacheIP)
	proxyChannel := listener.redisClient.BuildChannel(regPacket.ProxyGroupName, regPacket.ProxyServiceName, "*", types.RegisterServerData)

	regDataPacket := packets.RegisterServerData{
		Data:        packet.GetData(),
		AddressIP:   ip,
		Port:        port,
		ServiceName: regPacket.GetData().SenderServiceName,
		ServiceId:   regPacket.GetData().SenderServiceId,
	}

	cacheIP := packets.CacheIP{
		Data:      interfaces.NewPacketData(),
		AddressIP: ip,
		Port:      port,
	}

	go listener.redisClient.Publish(senderChannel, &cacheIP)
	go listener.redisClient.Publish(proxyChannel, &regDataPacket)

	return nil
}
