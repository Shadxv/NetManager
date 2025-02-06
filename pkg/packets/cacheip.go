package packets

import (
	"NetManager/pkg/interfaces"
	"NetManager/pkg/types"
)

type CacheIP struct {
	Data      interfaces.PacketData
	AddressIP string `json:"addressIP"`
	Port      int    `json:"port"`
}

func (packet *CacheIP) GetType() string {
	return types.CacheIP
}

func (packet *CacheIP) GetData() interfaces.PacketData {
	return packet.Data
}
