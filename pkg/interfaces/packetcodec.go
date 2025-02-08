package interfaces

type PacketCodec interface {
	UnmarshalPacket(packetType string, data []byte) (Packet, error)
	MarshalPacket(packet Packet) ([]byte, error)
	RegisterType(packetType string, creator func() Packet)
}
