package packets

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
)

type CacheIP struct {
	interfaces.PacketData
	AddressIP string `json:"addressIP"`
	Port      int    `json:"port"`
}

func (packet *CacheIP) GetType() string {
	return types.CacheIP
}
