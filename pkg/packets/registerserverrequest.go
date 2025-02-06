package packets

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
)

type RegisterServerRequest struct {
	Data             interfaces.PacketData
	ProxyGroupName   string `json:"proxyGroupName"`
	ProxyServiceName string `json:"proxyServiceName"`
}

func (packet *RegisterServerRequest) GetType() string {
	return types.RegisterServerRequest
}

func (packet *RegisterServerRequest) GetData() interfaces.PacketData {
	return packet.Data
}
