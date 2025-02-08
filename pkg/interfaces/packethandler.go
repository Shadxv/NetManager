package interfaces

type PacketHandler interface {
	GetType() string
	Handle(packet Packet) error
	GetRegistryType() Packet
}
