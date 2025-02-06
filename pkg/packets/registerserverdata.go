package packets

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
)

type RegisterServerData struct {
	Data        interfaces.PacketData
	AddressIP   string `json:"addressIP"`
	Port        int    `json:"port"`
	ServiceName string `json:"serviceName"`
	ServiceId   string `json:"serviceId"`
}

func (packet *RegisterServerData) GetType() string {
	return types.RegisterServerData
}

func (packet *RegisterServerData) GetData() interfaces.PacketData {
	return packet.Data
}
