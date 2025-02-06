package interfaces

type Packet interface {
	GetType() string
	GetData() PacketData
}
