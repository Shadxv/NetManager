package packets

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
)

type RegisterServerRequest struct {
	interfaces.PacketData
	ProxyGroupName   string `json:"proxyServiceGroup"`
	ProxyServiceName string `json:"proxyServiceName"`
}

func (packet *RegisterServerRequest) GetType() string {
	return types.RegisterServerRequest
}
