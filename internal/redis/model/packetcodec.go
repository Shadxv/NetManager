package model

import (
	"NetManager/pkg/interfaces"
	"encoding/json"
	"fmt"
)

type PacketCodec struct {
	typeRegistry map[string]func() interfaces.Packet
}

func NewPacketCodec() *PacketCodec {
	return &PacketCodec{
		typeRegistry: make(map[string]func() interfaces.Packet),
	}
}

func (codec *PacketCodec) UnmarshalPacket(packetType string, data []byte) (interfaces.Packet, error) {
	creator, ok := codec.typeRegistry[packetType]
	if !ok {
		return nil, fmt.Errorf("unknown packet type: %s", packetType)
	}

	packet := creator()
	if err := json.Unmarshal(data, packet); err != nil {
		return nil, err
	}

	return packet, nil
}

func (codec *PacketCodec) MarshalPacket(packet interfaces.Packet) ([]byte, error) {
	return json.Marshal(packet)
}

func (codec *PacketCodec) RegisterType(packetType string, creator func() interfaces.Packet) {
	codec.typeRegistry[packetType] = creator
}
